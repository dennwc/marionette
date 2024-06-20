package marionette

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"math"
	"strings"
	"time"
)

const (
	MARIONETTE_PROTOCOL_V3 = 3
	WEBDRIVER_ELEMENT_KEY  = "element-6066-11e4-a52e-4f735466cecf"
)

var RunningInDebugMode = false

type Client struct {
	SessionId    string
	Capabilities Capabilities

	tr *Transport
}

func NewClient() *Client {
	return &Client{
		tr: &Transport{},
	}
}

func (c *Client) Transport(tr *Transport) {
	c.tr = tr
}

func (c *Client) SessionID() string {
	return c.SessionId
}

func (c *Client) Connect(ctx context.Context, addr string) error {
	return c.tr.Connect(ctx, addr)
}

// NewSession create new session
func (c *Client) NewSession(sessionId string, cap *Capabilities) (*Response, error) {
	r, err := c.tr.Send("WebDriver:NewSession", map[string]any{
		"sessionId":    sessionId,
		"capabilities": cap,
	})
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(r.Value), &c)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteSession Marionette currently only accepts a session id, so if
// we call delete session can also close the TCP Connection
func (c *Client) DeleteSession() error {
	_, err := c.tr.Send("WebDriver:DeleteSession", nil)
	if err != nil {
		return err
	}
	return c.tr.Close()
}

// GetCapabilities informs the client of which WebDriver features are
// supported by Firefox and Marionette. They are immutable for the
// length of the session.
func (c *Client) GetCapabilities() (*Capabilities, error) {
	r, err := c.tr.Send("WebDriver:GetCapabilities", map[string]string{})
	if err != nil {
		return nil, err
	}

	var d = map[string]*Capabilities{}
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return nil, err
	}

	return d["capabilities"], nil
}

// SetScriptTimeout Set the timeout for asynchronous script execution.
func (c *Client) SetScriptTimeout(timeout time.Duration) (*Response, error) {
	data := map[string]int{"script": int(timeout.Milliseconds())}
	return c.SetTimeouts(data)
}

// SetImplicitTimout Set timeout for searching for elements.
func (c *Client) SetImplicitTimout(timeout time.Duration) (*Response, error) {
	data := map[string]int{"implicit": int(timeout.Milliseconds())}
	return c.SetTimeouts(data)
}

// SetPageLoadTimeout Set timeout for page loading.
func (c *Client) SetPageLoadTimeout(timeout time.Duration) (*Response, error) {
	data := map[string]int{"pageLoad": int(timeout.Milliseconds())}
	return c.SetTimeouts(data)
}

// SetTimeouts sets timeouts object.
//
// <dl>
// <dt><code>script</code> (number)
// <dd>Determines when to interrupt a script that is being evaluates.
//
// <dt><code>pageLoad</code> (number)
// <dd>Provides the timeout limit used to interrupt navigation of the
//
//	browsing context.
//
// <dt><code>implicit</code> (number)
// <dd>Gives the timeout of when to abort when locating an element.
// </dl>
func (c *Client) SetTimeouts(data map[string]int) (*Response, error) {
	r, err := c.tr.Send("WebDriver:SetTimeouts", data)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Get Timeouts Get current set timeouts
func (c *Client) GetTimeouts() (map[string]uint, error) {
	r, err := c.tr.Send("WebDriver:GetTimeouts", map[string]string{})
	if err != nil {
		return nil, err
	}

	var d = map[string]uint{}
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// Navigate open url
func (c *Client) Navigate(url string) (*Response, error) {
	r, err := c.tr.Send("WebDriver:Navigate", map[string]string{"url": url})
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Title get title
func (c *Client) Title() (string, error) {
	r, err := c.tr.Send("WebDriver:GetTitle", map[string]string{})
	if err != nil {
		return "", err
	}

	var d = map[string]string{}
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return "", err
	}

	return d["value"], nil
}

// URL get current url
func (c *Client) URL() (string, error) {
	var out struct {
		Value string `json:"value"`
	}
	err := c.tr.SendAndDecode(&out, "WebDriver:GetCurrentURL", nil)
	return out.Value, err
}

// Refresh the page.
func (c *Client) Refresh() error {
	_, err := c.tr.Send("WebDriver:Refresh", nil)
	return err
}

// Back go back in navigation history
func (c *Client) Back() error {
	_, err := c.tr.Send("WebDriver:Back", nil)
	return err
}

// Forward go forward in navigation history
func (c *Client) Forward() error {
	_, err := c.tr.Send("WebDriver:Forward", nil)
	return err
}

// SetContext Sets the context of the subsequent commands to be either "chrome" or "content".
// Must be one of "chrome" or "content" only.
func (c *Client) SetContext(value Context) (*Response, error) {
	return c.tr.Send("Marionette:SetContext", map[string]string{"value": fmt.Sprint(value)})
}

// Context Gets the context of the server, either "chrome" or "content".
func (c *Client) Context() (*Response, error) {
	return c.tr.Send("Marionette:GetContext", nil)
}

// GetWindowHandle returns the current window ID
func (c *Client) GetWindowHandle() (string, error) {
	r, err := c.tr.Send("WebDriver:GetWindowHandle", nil)
	if err != nil {
		return "", err
	}

	var d map[string]string
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return "", err
	}
	return d["value"], nil
}

// GetWindowHandles return array of window ID currently opened
func (c *Client) GetWindowHandles() ([]string, error) {
	r, err := c.tr.Send("WebDriver:GetWindowHandles", nil)
	if err != nil {
		return nil, err
	}

	var d []string
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// SwitchToWindow switch to specific window.
func (c *Client) SwitchToWindow(name string) error {
	_, err := c.tr.Send("WebDriver:SwitchToWindow", map[string]any{"focus": true, "handle": name})
	return err
}

// GetWindowRect gets window position and size
func (c *Client) GetWindowRect() (rect *WindowRect, err error) {
	r, err := c.tr.Send("WebDriver:GetWindowRect", nil)
	if err != nil {
		return nil, err
	}

	rect = new(WindowRect)
	err = json.Unmarshal([]byte(r.Value), &rect)
	if err != nil {
		return nil, err
	}

	return
}

// SetWindowRect sets window position and size
func (c *Client) SetWindowRect(rect WindowRect) error {
	_, err := c.tr.Send("WebDriver:SetWindowRect", map[string]any{
		"x":      rect.X,
		"y":      rect.Y,
		"width":  math.Floor(rect.Width),
		"height": math.Floor(rect.Height),
	})
	return err
}

// MaximizeWindow maximizes window.
func (c *Client) MaximizeWindow() (*WindowRect, error) {
	rect := new(WindowRect)
	err := c.tr.SendAndDecode(rect, "WebDriver:MaximizeWindow", nil)
	if err != nil {
		return nil, err
	}
	return rect, nil
}

// MinimizeWindow Synchronously minimizes the user agent window as if the user pressed
// the minimize button.
func (c *Client) MinimizeWindow() (*WindowRect, error) {
	rect := new(WindowRect)
	err := c.tr.SendAndDecode(rect, "WebDriver:MinimizeWindow", nil)
	if err != nil {
		return nil, err
	}
	return rect, nil
}

// FullscreenWindow Synchronously sets the user agent window to full screen as if the user
// had done "View > Enter Full Screen"
func (c *Client) FullscreenWindow() (rect *WindowRect, err error) {
	r, err := c.tr.Send("WebDriver:FullscreenWindow", nil)
	if err != nil {
		return nil, err
	}

	rect = new(WindowRect)
	err = json.Unmarshal([]byte(r.Value), &rect)
	if err != nil {
		return nil, err
	}

	return
}

// NewWindow opens a new top-level browsing context window.
//
// param: type string
// Optional type of the new top-level browsing context. Can be one of
// `tab` or `window`. Defaults to `tab`.
//
// param: focus bool
// Optional flag if the new top-level browsing context should be opened
// in foreground (focused) or background (not focused). Defaults to false.
//
// param: private bool
// Optional flag, which gets only evaluated for type `window`. True if the
// new top-level browsing context should be a private window.
// Defaults to false.
//
// return {"handle": string, "type": string}
// Handle and type of the new browsing context.
func (c *Client) NewWindow(focus bool, typ string, private bool) (*Response, error) {
	//TODO: would be nice if we could create a Window struct and return that struct instead of the Response object
	return c.tr.Send("WebDriver:NewWindow", map[string]any{
		"focus":   focus,
		"type":    typ,
		"private": private,
	})
}

// CloseWindow closes current window.
func (c *Client) CloseWindow() (*Response, error) {
	return c.tr.Send("WebDriver:CloseWindow", nil)
}

// CloseChromeWindow closes the currently selected chrome window.
//
// If it is the last window currently open, the chrome window will not be
// closed to prevent a shutdown of Firefox. Instead the returned
// list of chrome window handles is empty.
//
// return []string
// Unique chrome window handles of remaining chrome windows.
//
// error NoSuchWindowError
// Top-level browsing context has been discarded.
func (c *Client) CloseChromeWindow() (*Response, error) {
	return c.tr.Send("WebDriver:CloseChromeWindow", nil)
}

// SwitchToFrame switch to frame - strategies: By(ID), By(NAME) or name only.
func (c *Client) SwitchToFrame(by By, value string) error {
	//with current marionette implementation we have to find the element first and send the switchToFrame
	//command with the UUID, else it wont work.
	//https://bugzilla.mozilla.org/show_bug.cgi?id=1143908
	frame, err := c.FindElement(by, value)
	if err != nil {
		return err
	}

	_, err = c.tr.Send("WebDriver:SwitchToFrame", map[string]any{"element": frame.Id(), "focus": true})
	return err
}

// SwitchToParentFrame switch to parent frame
func (c *Client) SwitchToParentFrame() error {
	_, err := c.tr.Send("WebDriver:SwitchToParentFrame", nil)
	return err
}

// AddCookie Adds a cookie
func (c *Client) AddCookie(cookie Cookie) (*Response, error) {
	return c.tr.Send("WebDriver:AddCookie", map[string]any{"cookie": cookie})
}

// GetCookies Get all cookies
func (c *Client) GetCookies() ([]Cookie, error) {
	r, err := c.tr.Send("WebDriver:GetCookies", nil)
	if err != nil {
		return nil, err
	}

	var cookies []Cookie
	_ = json.Unmarshal([]byte(r.Value), &cookies)

	return cookies, nil
}

// DeleteCookie Deletes cookie by name
func (c *Client) DeleteCookie(name string) (error, error) {
	_, err := c.tr.Send("WebDriver:DeleteCookie", map[string]any{"name": name})
	return err, nil
}

// DeleteAllCookies Delete all cookies
func (c *Client) DeleteAllCookies() error {
	_, err := c.tr.Send("WebDriver:DeleteAllCookies", nil)
	return err
}

//////////////////
// WEB ELEMENTS //
//////////////////

func isElementEnabled(c *Client, id string) bool {
	r, err := c.tr.Send("WebDriver:IsElementEnabled", map[string]any{"id": id})
	if err != nil {
		return false
	}

	return strings.Contains(r.Value, "\"value\":true")
}

func isElementSelected(c *Client, id string) bool {
	r, err := c.tr.Send("WebDriver:IsElementSelected", map[string]any{"id": id})
	if err != nil {
		return false
	}

	return strings.Contains(r.Value, "\"value\":true")
}

func isElementDisplayed(c *Client, id string) bool {
	r, err := c.tr.Send("WebDriver:IsElementDisplayed", map[string]any{"id": id})
	if err != nil {
		return false
	}

	return strings.Contains(r.Value, "\"value\":true")
}

func getElementTagName(c *Client, id string) string {
	r, err := c.tr.Send("WebDriver:GetElementTagName", map[string]any{"id": id})
	if err != nil {
		return ""
	}

	var d = map[string]string{}
	json.Unmarshal([]byte(r.Value), &d)

	return d["value"]
}

func getElementText(c *Client, id string) string {
	r, err := c.tr.Send("WebDriver:GetElementText", map[string]any{"id": id})
	if err != nil {
		return ""
	}

	var d = map[string]string{}
	json.Unmarshal([]byte(r.Value), &d)

	return d["value"]
}

func (c *Client) findElements(by By, value string, startNode *string) ([]*WebElement, error) {
	var params map[string]any
	if startNode == nil || *startNode == "" {
		params = map[string]any{"using": fmt.Sprint(by), "value": value}
	} else {
		params = map[string]any{"using": fmt.Sprint(by), "value": value, "element": *startNode}
	}

	r, err := c.tr.Send("WebDriver:FindElements", params)
	if err != nil {
		return nil, err
	}

	var d []map[string]string
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return nil, err
	}

	var e []*WebElement
	for _, v := range d {
		e = append(e, &WebElement{c: c, id: v[WEBDRIVER_ELEMENT_KEY]})
	}
	return e, nil
}

// FindElements Find elements using the indicated search strategy.
func (c *Client) FindElements(by By, value string) ([]*WebElement, error) {
	return c.findElements(by, value, nil)
}

func (c *Client) findElement(by By, value string, startNode *string) (*WebElement, error) {
	var params map[string]string
	if startNode == nil || *startNode == "" {
		params = map[string]string{"using": fmt.Sprint(by), "value": value}
	} else {
		params = map[string]string{"using": fmt.Sprint(by), "value": value, "element": *startNode}
	}
	e := &WebElement{c: c}
	err := c.tr.SendAndDecode(e, "WebDriver:FindElement", params)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// FindElement Find an element using the indicated search strategy.
func (c *Client) FindElement(by By, value string) (*WebElement, error) {
	return c.findElement(by, value, nil)
}

// GetActiveElement Returns the page's active element.
func (c *Client) GetActiveElement() (*WebElement, error) {
	e := &WebElement{c: c}
	err := c.tr.SendAndDecode(e, "WebDriver:GetActiveElement", nil)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// PageSource get page source
func (c *Client) PageSource() (string, error) {
	var out struct {
		Value string `json:"value"`
	}
	err := c.tr.SendAndDecode(&out, "WebDriver:GetPageSource", nil)
	return out.Value, err
}

func convertScriptArgs(args []any) {
	for i, arg := range args {
		if e, ok := arg.(*WebElement); ok {
			args[i] = e.Id() // TODO: convert properly
		}
	}
}

// ExecuteScript Execute JS Script
func (c *Client) ExecuteScript(script string, args []any, timeout time.Duration, newSandbox bool) (*Response, error) {
	convertScriptArgs(args)
	return c.tr.Send("WebDriver:ExecuteScript", map[string]any{
		"scriptTimeout": int(timeout.Milliseconds()),
		"script":        script,
		"args":          args,
		"newSandbox":    newSandbox,
	})
}

// ExecuteAsyncScript Execute JS Script Async
// TODO: Add missing arguments/options
func (c *Client) ExecuteAsyncScript(script string, args []any, newSandbox bool) (*Response, error) {
	convertScriptArgs(args)
	return c.tr.Send("WebDriver:ExecuteAsyncScript", map[string]any{
		"script":     script,
		"args":       args,
		"newSandbox": newSandbox,
	})
}

// DismissAlert dismisses the dialog - like clicking No/Cancel
func (c *Client) DismissAlert() error {
	_, err := c.tr.Send("WebDriver:DismissAlert", nil)
	return err
}

// AcceptAlert accepts the dialog - like clicking Ok/Yes
func (c *Client) AcceptAlert() error {
	_, err := c.tr.Send("WebDriver:AcceptAlert", nil)
	return err
}

// TextFromAlert gets text from the dialog
func (c *Client) TextFromAlert() (string, error) {
	r, err := c.tr.Send("WebDriver:GetAlertText", map[string]any{"key": "value"})
	if err != nil {
		return "", err
	}

	var d = map[string]string{}
	json.Unmarshal([]byte(r.Value), &d)

	return d["value"], nil
}

// SendAlertText sends text to a dialog
func (c *Client) SendAlertText(keys string) error {
	_, err := c.tr.Send("WebDriver:SendAlertText", map[string]any{"text": keys})
	return err
}

// Quit quits the session and request browser process to terminate.
func (c *Client) Quit() (*Response, error) {
	return c.tr.Send("Marionette:Quit", map[string][]string{"flags": {"eForceQuit"}})
}

func (c *Client) takeScreenshot(startNode *string) ([]byte, error) {
	var params map[string]string
	if startNode == nil || *startNode == "" {
		params = map[string]string{}
	} else {
		params = map[string]string{"id": *startNode}
	}
	var out struct {
		Value []byte `json:"value"`
	}
	err := c.tr.SendAndDecode(&out, "WebDriver:TakeScreenshot", params)
	return out.Value, err
}

func (c *Client) takeScreenshotImage(startNode *string) (image.Image, error) {
	data, err := c.takeScreenshot(startNode)
	if err != nil {
		return nil, err
	}
	return png.Decode(bytes.NewReader(data))
}

// Screenshot takes a screenshot of the page.
func (c *Client) Screenshot() ([]byte, error) {
	return c.takeScreenshot(nil)
}

// ScreenshotImage takes a screenshot of the page.
func (c *Client) ScreenshotImage() (image.Image, error) {
	return c.takeScreenshotImage(nil)
}

func (c *Client) PerformActions(actions Actions) (*Response, error) {
	r, err := c.tr.Send("WebDriver:PerformActions", actions)
	return r, err
}
