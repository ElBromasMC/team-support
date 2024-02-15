// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package layout

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Header() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<header class=\"sticky top-0 z-50 text-slate-600 border-b border-black bg-white\"><div class=\"flex items-center px-4 py-4 max-w-7xl mx-auto lg:py-9\"><!-- Navigation button --><script>\n                function handleNavbarDisplay(el) {\n                    if (el.dataset.open != null) {\n                        delete el.dataset.open\n                    } else {\n                        el.dataset.open = \"\"\n                    }\n                }\n            </script><button class=\"group peer flex items-center justify-center w-8 h-8 lg:hidden\" onclick=\"handleNavbarDisplay(this)\" type=\"button\"><svg class=\"w-4 h-4 group-data-[open]:hidden\" aria-hidden=\"true\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"currentColor\" viewBox=\"0 0 4 15\"><path d=\"M3.5 1.5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Zm0 6.041a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Zm0 5.959a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z\"></path></svg> <svg class=\"hidden w-5 h-5 group-data-[open]:block\" aria-hidden=\"true\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\"><path stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18 18 6m0 12L6 6\"></path></svg></button><!-- Navigation bar (mobiles) --><nav class=\"hidden absolute top-full inset-x-0 bg-white peer-data-[open]:block lg:!hidden\"><ul class=\"flex flex-col divide-y border-y border-black\"><li><a class=\"flex justify-center items-center h-12\" href=\"#\">Nosotros</a></li><li><a class=\"flex justify-center items-center h-12\" href=\"#\">Servicio</a></li><li><a class=\"flex justify-center items-center h-12\" href=\"#\">Contáctenos</a></li><li><a class=\"flex justify-center items-center h-12\" href=\"/store\">Tienda</a></li><li><a class=\"flex justify-center items-center h-12\" href=\"#\">Garantía</a></li><li><a class=\"flex justify-center items-center h-12\" href=\"#\">Intranet</a></li></ul></nav><!-- Logo --><a class=\"flex items-center justify-center h-8 ml-3 lg:ml-0\" href=\"/\"><img class=\"w-44 lg:w-[210px]\" src=\"/static/img/logo1.png\" alt=\"Team Support Services\"></a><!-- Navigation bar --><nav class=\"hidden grow ml-16 font-semibold text-sm lg:block\"><ul class=\"flex justify-between\"><li><a class=\"hover:text-slate-900\" href=\"#\">Nosotros</a></li><li><a class=\"hover:text-slate-900\" href=\"#\">Servicio</a></li><li><a class=\"hover:text-slate-900\" href=\"#\">Contáctenos</a></li><li><a class=\"hover:text-slate-900\" href=\"/store\">Tienda</a></li><li><a class=\"hover:text-slate-900\" href=\"#\">Garantía</a></li><li><a class=\"hover:text-slate-900\" href=\"#\">Intranet</a></li></ul></nav><!-- Shopping cart --><a class=\"flex items-center justify-center w-8 h-8 ml-auto hover:text-slate-900 lg:ml-16\" href=\"#\"><svg class=\"w-6 h-6 lg:w-8 lg:h-8\" aria-hidden=\"true\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 18 20\"><path stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"1.5\" d=\"M6 15a2 2 0 1 0 0 4 2 2 0 0 0 0-4Zm0 0h8m-8 0-1-4m9 4a2 2 0 1 0 0 4 2 2 0 0 0 0-4Zm-9-4h10l2-7H3m2 7L3 4m0 0-.792-3H1\"></path></svg></a><!-- User profile --><a class=\"flex items-center justify-center w-8 h-8 ml-3 hover:text-slate-900 lg:ml-16\" href=\"/signup\"><svg class=\"w-6 h-6 lg:w-8 lg:h-8\" aria-hidden=\"true\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"currentColor\" viewBox=\"0 0 20 20\"><path d=\"M10 0a10 10 0 1 0 10 10A10.011 10.011 0 0 0 10 0Zm0 5a3 3 0 1 1 0 6 3 3 0 0 1 0-6Zm0 13a8.949 8.949 0 0 1-4.951-1.488A3.987 3.987 0 0 1 9 13h2a3.987 3.987 0 0 1 3.951 3.512A8.949 8.949 0 0 1 10 18Z\"></path></svg></a></div></header>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
