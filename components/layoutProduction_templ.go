// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.731
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Base(children ...templ.Component) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><title>Dashboard</title><link href=\"/static/css/tailwind.css\" rel=\"stylesheet\"><script src=\"https://unpkg.com/htmx.org@2.0.0\"></script><link rel=\"stylesheet\" href=\"https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@24,400,0,0\"><link rel=\"preconnect\" href=\"https://fonts.googleapis.com\"><link rel=\"preconnect\" href=\"https://fonts.gstatic.com\" crossorigin><link href=\"https://fonts.googleapis.com/css2?family=Roboto:ital,wght@0,100;0,300;0,400;0,500;0,700;0,900;1,100;1,300;1,400;1,500;1,700;1,900&amp;display=swap\" rel=\"stylesheet\"><style>\r\n        body.dashboard #dash_under{\r\n            text-decoration: underline red 2px ;\r\n            text-underline-offset: 5px;\r\n\r\n        }\r\n        body.tasks #tasks_under{\r\n            text-decoration: underline red 2px ;\r\n            text-underline-offset: 5px;\r\n        }\r\n        body.pomodoro #pomodoro_under{\r\n            text-decoration: underline red 2px ;\r\n            text-underline-offset: 5px;\r\n        }\r\n\r\n\r\n\r\n    </style></head><body class=\"font-roboto text-xl\"><div id=\"nav-container\" class=\"drop-shadow-md border-b border-shadow bg-white-900 m-0 p-0\"><header class=\"box-border justify-between flex items-center\"><nav class=\"w-full\"><ul class=\"flex justify-between items-center p-3\"><li class=\"inline-block p-4\"><a href=\"#\">HOME</a></li><div class=\"flex\"><li id=\"dash_under\" class=\"inline-block p-4\"><a href=\"/api/dashboard\" class=\"decoration-solid\">DASHBOARD</a></li><li id=\"tasks_under\" class=\"inline-block p-4\"><a href=\"#\">TASKS</a></li><li id=\"pomodoro_under\" class=\"inline-block p-4\"><a href=\"#\">POMODORO</a></li></div><li class=\"inline-block p-4\"><a href=\"#\">PROFILE</a></li></ul></nav></header></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, child := range children {
			templ_7745c5c3_Err = child.Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script>\r\n    document.addEventListener(\"DOMContentLoaded\", function() {\r\n        var bodyClass = \"\";\r\n        var currentPage = window.location.pathname;\r\n        if (currentPage === \"/api/test/dashboard\") {\r\n            bodyClass = \"dashboard\";\r\n        }\r\n        if (currentPage === \"/api/test/tasks\") {\r\n            bodyClass = \"tasks\";\r\n        }\r\n        if (currentPage === \"/api/test/pomodoro\") {\r\n            bodyClass = \"pomodoro\";\r\n        }\r\n        if (bodyClass) {\r\n            document.body.classList.add(bodyClass);\r\n        }\r\n    });\r\n</script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}