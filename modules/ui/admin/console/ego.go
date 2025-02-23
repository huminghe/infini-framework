// Generated by ego.
// DO NOT EDIT

package console

import (
	"fmt"
	"github.com/huminghe/infini-framework/modules/ui/common"
	"io"
	"net/http"
)

var _ = fmt.Sprint("") // just so that we can keep the fmt import for now
func Index(w http.ResponseWriter, r *http.Request) error {
	_, _ = io.WriteString(w, "\n\n")
	_, _ = io.WriteString(w, "\n")
	_, _ = io.WriteString(w, "\n\n")
	common.Head(w, "Console", "")
	_, _ = io.WriteString(w, "\n\n<link rel=\"stylesheet\" href=\"/static/assets/uikit-2.27.1/css/components/notify.min.css\" />\n<script src=\"/static/assets/js/jquery.hotkeys-0.7.9.min.js\"></script>\n<script src=\"/static/assets/uikit-2.27.1/js/components/notify.min.js\"></script>\n<script src=\"/static/assets/js/page/console.js?v=1\"></script>\n<style type=\"text/css\">\n    pre{\n        margin-top:0px;\n        margin-bottom:0px;\n        padding-top:4px;\n        padding-bottom:4px;\n    }\n    .request pre{\n        color:royalblue;\n    }\n    #log {\n        max-height: 200px;\n    }\n    #logging {\n        height: 500px;\n        max-height: 1000px;\n    }\n</style>\n")
	common.Body(w)
	_, _ = io.WriteString(w, "\n")
	common.Nav(w, r, "Console")
	_, _ = io.WriteString(w, "\n\n\n\n<div class=\"tm-middle\">\n\n    <div class=\"uk-container uk-container-center\">\n\n        <div class=\"uk-grid \" data-uk-grid-margin>\n            <div class=\"uk-width-1-1 \">\n\n                <div class=\"uk-alert uk-alert-success\" data-uk-alert=\"\">\n                    <a href=\"#\" class=\"uk-alert-close uk-close\"></a>\n                    <p id=\"connect_status\"></p>\n                </div>\n\n                <div id=\"log\" class=\"uk-scrollable-text\"></div>\n                <form id=\"form\" class=\"uk-form\">\n                    <input id=\"msg\" type=\"text\"  style=\"ime-mode: disabled\" title=\"Press <s> to quick input; Type HELP for more details.\" placeholder=\"Press <s> to quick input; Type HELP for more details.\" autocomplete=\"on\" class=\"uk-width-1-1\">\n                </form>\n            </div>\n            <div class=\"uk-width-1-1 \">\n\n                <form class=\"uk-form\" onsubmit=\"return false;\">\n\n                    <fieldset data-uk-margin>\n                        <legend class=\"uk-text-primary\">Realtime logging</legend>\n                        <select id=\"log_level\">\n                            <option>ERROR</option>\n                            <option>WARN</option>\n                            <option selected>INFO</option>\n                            <option>DEBUG</option>\n                            <option>TRACE</option>\n                        </select>\n                        <input id=\"file_pattern\" type=\"text\" class=\"uk-form-width-medium\" placeholder=\"FilePattern, eg: crawler*.go\">\n                        <input id=\"func_pattern\" type=\"text\" class=\"uk-form-width-medium\" placeholder=\"FuncPattern, eg: runPipeline*\">\n                        <input id=\"message_pattern\" type=\"text\" class=\"uk-form-width-medium\" placeholder=\"MessagePattern, eg: *timeout*\">\n                        <button class=\"uk-button uk-icon-play uk-button-success\" onclick=\"updateLogging(true)\"> START</button>\n                        <button class=\"uk-button uk-icon-stop uk-button-danger\" disabled onclick=\"updateLogging(false)\"> STOP</button>\n                        <div id=\"update_config_response\"></div>\n                    </fieldset>\n                </form>\n\n\n                <div id=\"logging\" class=\"uk-scrollable-text \"></div>\n            </div>\n        </div>\n\n    </div>\n\n</div>\n\n")
	common.Footer(w)
	_, _ = io.WriteString(w, "\n")
	return nil
}
