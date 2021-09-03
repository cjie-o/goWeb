package weblib

func (B *WebBase) WebErr(msg string) {
	B.IsErr().SetBody(msg).Send()
}
func (B *WebBase) WebSucess(msg string) {
	B.SetBody(msg).Send()
}
func (B *WebBase) Web403() {
	B.SetCode(403).SetBody("非法请求，已记录").IsErr().Send()
}
func (B *WebBase) Web404() {
	B.SetCode(404).SetBody("资源不存在").IsErr().Send()
}
func (B *WebBase) Web500() {
	B.SetCode(500).SetBody("服务器异常").IsErr().Send()
}
func (B *WebBase) WebPage(page string) {
	B.SetBody(page).IsJson(false).Send()
}

// func (B *WebBase) WebCtype() {}
func (B *WebBase) WebLocation(url string) {
	B.SetHead("location", url).SetCode(302).Send()
}
func (B *WebBase) WebJson(data string) {
	B.SetBody(data).Send()
}
