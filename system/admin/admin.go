// Package admin desrcibes the admin view containing references to
// various managers and editors
package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/extvos/ponzu/system/admin/user"
	"github.com/extvos/ponzu/system/api/analytics"
	"github.com/extvos/ponzu/system/db"
	"github.com/extvos/ponzu/system/item"
)

var startAdminHTML = `<!doctype html>
<html lang="en">
    <head>
        <title>{{ .Logo }}</title>
        <script type="text/javascript" src="/admin/static/common/js/jquery-2.1.4.min.js"></script>
        <script type="text/javascript" src="/admin/static/common/js/util.js"></script>
        <script type="text/javascript" src="/admin/static/dashboard/js/materialize.min.js"></script>
        <script type="text/javascript" src="/admin/static/dashboard/js/chart.bundle.min.js"></script>
        <script type="text/javascript" src="/admin/static/editor/js/materialNote.js"></script> 
        <script type="text/javascript" src="/admin/static/editor/js/ckMaterializeOverrides.js"></script>
                  
        <link rel="stylesheet" href="/admin/static/dashboard/css/material-icons.css" />     
        <link rel="stylesheet" href="/admin/static/dashboard/css/materialize.min.css" />
        <link rel="stylesheet" href="/admin/static/editor/css/materialNote.css" />
        <link rel="stylesheet" href="/admin/static/dashboard/css/admin.css" />    

        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
    </head>
    <body class="grey lighten-4">
       <div class="navbar-fixed">
            <nav class="grey darken-2">
            <div class="nav-wrapper">
                <a class="brand-logo" href="/admin">{{ .Logo }}</a>

                <ul class="right">
                    <li><a href="/admin/logout">注销</a></li>
                </ul>
            </div>
            </nav>
        </div>

        <div class="admin-ui row">`

var mainAdminHTML = `
            <div class="left-nav col s3">
                <div class="card">
                <ul class="card-content collection">
                    <div class="card-title">内容</div>
                                    
                    {{ range $t, $f := .Types }}
                    <div class="row collection-item">
                        <li><a class="col s12" href="/admin/contents?type={{ $t }}"><i class="tiny left material-icons">playlist_add</i>{{ $f }}</a></li>
                    </div>
                    {{ end }}

                    <div class="card-title">系统</div>                                
                    <div class="row collection-item">
                        <li><a class="col s12" href="/admin/configure"><i class="tiny left material-icons">settings</i>设置</a></li>
                        <li><a class="col s12" href="/admin/configure/users"><i class="tiny left material-icons">supervisor_account</i>管理员</a></li>
                        <li><a class="col s12" href="/admin/uploads"><i class="tiny left material-icons">swap_vert</i>上传</a></li>
                        <li><a class="col s12" href="/admin/addons"><i class="tiny left material-icons">settings_input_svideo</i>组件</a></li>
                    </div>
                </ul>
                </div>
            </div>
            {{ if .Subview}}
            <div class="subview col s9">
                {{ .Subview }}
            </div>
            {{ end }}`

var endAdminHTML = `
        </div>
        <footer class="row">
            <div class="col s12">
                <p class="center-align">Powered by &copy; Expeak  &nbsp;&vert;&nbsp;</p>
            </div>     
        </footer>
    </body>
</html>`

type admin struct {
	Logo    string
	Types   map[string]string
	Subview template.HTML
}

// Admin ...
func Admin(view []byte) (_ []byte, err error) {
	cfg, err := db.Config("name")
	if err != nil {
		return
	}

	if cfg == nil {
		cfg = []byte("")
	}

	types := make(map[string]string)
	for k, f := range item.Types {
		types[k] = k
		if vv, ok := f().(item.Identifiable); ok {
			types[k] = vv.TypeName()
		}
	}

	a := admin{
		Logo:    string(cfg),
		Types:   types, //item.Types,
		Subview: template.HTML(view),
	}

	buf := &bytes.Buffer{}
	html := startAdminHTML + mainAdminHTML + endAdminHTML
	tmpl := template.Must(template.New("admin").Parse(html))
	err = tmpl.Execute(buf, a)
	if err != nil {
		return
	}

	return buf.Bytes(), nil
}

var initAdminHTML = `
<div class="init col s5">
<div class="card">
<div class="card-content">
    <div class="card-title">欢迎！</div>
    <blockquote>在您使用系统之前需要填写下面的配置以初始化。所有的信息后面还可以更新，但是当前您需要先初始化才能开始使用。</blockquote>
    <form method="post" action="/admin/init" class="row">
        <div>设置</div>
        <div class="input-field col s12">        
            <input placeholder="输入您的站点名称（仅内部使用）" class="validate required" type="text" id="name" name="name"/>
            <label for="name" class="active">站点名称</label>
        </div>
        <div class="input-field col s12">        
            <input placeholder="获取SSL证书时候需要（例如： www.example.com 或  example.com）" class="validate" type="text" id="domain" name="domain"/>
            <label for="domain" class="active">域名</label>
        </div>
        <div>管理员</div>
        <div class="input-field col s12">
            <input placeholder="您的邮箱地址，比如：you@example.com" class="validate required" type="email" id="email" name="email"/>
            <label for="email" class="active">邮箱</label>
        </div>
        <div class="input-field col s12">
            <input placeholder="搞个厉害点的密码吧" class="validate required" type="password" id="password" name="password"/>
            <label for="password" class="active">密码</label>        
        </div>
        <button class="btn waves-effect waves-light right">开始</button>
    </form>
</div>
</div>
</div>
<script>
    $(function() {
        $('.nav-wrapper ul.right').hide();
        
        var logo = $('a.brand-logo');
        var name = $('input#name');
        var domain = $('input#domain');
        var hostname = domain.val();

        if (hostname === '') {    
            hostname = window.location.host || window.location.hostname;
        }
        
        if (hostname.indexOf(':') !== -1) {
            hostname = hostname.split(':')[0];
        }
        
        domain.val(hostname);
        
        name.on('change', function(e) {
            logo.text(e.target.value);
        });

    });
</script>
`

// Init ...
func Init() ([]byte, error) {
	html := startAdminHTML + initAdminHTML + endAdminHTML

	name, err := db.Config("name")
	if err != nil {
		return nil, err
	}

	if name == nil {
		name = []byte("")
	}

	a := admin{
		Logo: string(name),
	}

	buf := &bytes.Buffer{}
	tmpl := template.Must(template.New("init").Parse(html))
	err = tmpl.Execute(buf, a)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var loginAdminHTML = `
<div class="init col s5">
<div class="card">
<div class="card-content">
    <div class="card-title">欢迎！</div>
    <blockquote>请使用您的邮箱地址和密码登录系统！</blockquote>
    <form method="post" action="/admin/login" class="row">
        <div class="input-field col s12">
            <input placeholder="您的邮箱地址，比如：you@example.com" class="validate required" type="email" id="email" name="email"/>
            <label for="email" class="active">邮箱</label>
        </div>
        <div class="input-field col s12">
            <input placeholder="输入密码" class="validate required" type="password" id="password" name="password"/>
            <a href="/admin/recover">密码忘了？</a>            
            <label for="password" class="active">密码</label>  
        </div>
        <button class="btn waves-effect waves-light right">登录</button>
    </form>
</div>
</div>
</div>
<script>
    $(function() {
        $('.nav-wrapper ul.right').hide();
    });
</script>
`

// Login ...
func Login() ([]byte, error) {
	html := startAdminHTML + loginAdminHTML + endAdminHTML

	cfg, err := db.Config("name")
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		cfg = []byte("")
	}

	a := admin{
		Logo: string(cfg),
	}

	buf := &bytes.Buffer{}
	tmpl := template.Must(template.New("login").Parse(html))
	err = tmpl.Execute(buf, a)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var forgotPasswordHTML = `
<div class="init col s5">
<div class="card">
<div class="card-content">
    <div class="card-title">账户重置</div>
    <blockquote>请输入您的帐户邮箱地址，您这个邮箱地址将会收到一封重置账户的邮件，然后根据邮件指引重置您的帐户。别忘了检查一下邮箱的垃圾信件分组哦！</blockquote>
    <form method="post" action="/admin/recover" class="row" enctype="multipart/form-data">
        <div class="input-field col s12">
            <input placeholder="您的邮箱地址，比如：you@example.com" class="validate required" type="email" id="email" name="email"/>
            <label for="email" class="active">邮箱</label>
        </div>
        
        <a href="/admin/recover/key">已经有一个重置验证码？</a>
        <button class="btn waves-effect waves-light right">发送重置邮件</button>
    </form>
</div>
</div>
</div>
<script>
    $(function() {
        $('.nav-wrapper ul.right').hide();
    });
</script>
`

// ForgotPassword ...
func ForgotPassword() ([]byte, error) {
	html := startAdminHTML + forgotPasswordHTML + endAdminHTML

	cfg, err := db.Config("name")
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		cfg = []byte("")
	}

	a := admin{
		Logo: string(cfg),
	}

	buf := &bytes.Buffer{}
	tmpl := template.Must(template.New("forgotPassword").Parse(html))
	err = tmpl.Execute(buf, a)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var recoveryKeyHTML = `
<div class="init col s5">
<div class="card">
<div class="card-content">
    <div class="card-title">账户重置</div>
    <blockquote>请检查一下您提供的邮箱地址是否收到一封包含了一个重置验证码的邮件，别忘了垃圾信件分组也找找！</blockquote>
    <form method="post" action="/admin/recover/key" class="row" enctype="multipart/form-data">
        <div class="input-field col s12">
            <input placeholder="输入重置验证码" class="validate required" type="text" id="key" name="key"/>
            <label for="key" class="active">重置验证码</label>
        </div>

        <div class="input-field col s12">
            <input placeholder="您的邮箱地址，比如：you@example.com" class="validate required" type="email" id="email" name="email"/>
            <label for="email" class="active">邮箱地址</label>
        </div>

        <div class="input-field col s12">
            <input placeholder="输入新密码" class="validate required" type="password" id="password" name="password"/>
            <label for="password" class="active">新密码</label>
        </div>
        
        <button class="btn waves-effect waves-light right">更新账户</button>
    </form>
</div>
</div>
</div>
<script>
    $(function() {
        $('.nav-wrapper ul.right').hide();
    });
</script>
`

// RecoveryKey ...
func RecoveryKey() ([]byte, error) {
	html := startAdminHTML + recoveryKeyHTML + endAdminHTML

	cfg, err := db.Config("name")
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		cfg = []byte("")
	}

	a := admin{
		Logo: string(cfg),
	}

	buf := &bytes.Buffer{}
	tmpl := template.Must(template.New("recoveryKey").Parse(html))
	err = tmpl.Execute(buf, a)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UsersList ...
func UsersList(req *http.Request) ([]byte, error) {
	html := `
    <div class="card user-management">
        <div class="card-title">编辑您的帐户</div>    
        <form class="row" enctype="multipart/form-data" action="/admin/configure/users/edit" method="post">
            <div class="col s9">
                <label class="active">邮件地址</label>
                <input type="email" name="email" value="{{ .User.Email }}"/>
            </div>

            <div class="col s9">
                <div>请输入您的密码以继续更改</div>
                
                <label class="active">当前密码</label>
                <input type="password" name="password"/>
            </div>

            <div class="col s9">
                <label class="active">新密码：（留空将不会更新密码）</label>
                <input name="new_password" type="password"/>
            </div>

            <div class="col s9">                        
                <button class="btn waves-effect waves-light green right" type="submit">保存</button>
            </div>
        </form>

        <div class="card-title">添加一个新用户</div>        
        <form class="row" enctype="multipart/form-data" action="/admin/configure/users" method="post">
            <div class="col s9">
                <label class="active">邮件地址</label>
                <input type="email" name="email" value=""/>
            </div>

            <div class="col s9">
                <label class="active">密码</label>
                <input type="password" name="password"/>
            </div>

            <div class="col s9">            
                <button class="btn waves-effect waves-light green right" type="submit">添加用户</button>
            </div>   
        </form>        

        <div class="card-title">删除管理员用户</div>        
        <ul class="users row">
            {{ range .Users }}
            <li class="col s9">
                {{ .Email }}
                <form enctype="multipart/form-data" class="delete-user __ponzu right" action="/admin/configure/users/delete" method="post">
                    <span>删除</span>
                    <input type="hidden" name="email" value="{{ .Email }}"/>
                    <input type="hidden" name="id" value="{{ .ID }}"/>
                </form>
            </li>
            {{ end }}
        </ul>
    </div>
    `
	script := `
    <script>
        $(function() {
            var del = $('.delete-user.__ponzu span');
            del.on('click', function(e) {
                if (confirm("请确认：\n\n您是否确定删除该用户？\n本操作无法恢复！")) {
                    $(e.target).parent().submit();
                }
            });
        });
    </script>
    `
	// get current user out to pass as data to execute template
	j, err := db.CurrentUser(req)
	if err != nil {
		return nil, err
	}

	var usr user.User
	err = json.Unmarshal(j, &usr)
	if err != nil {
		return nil, err
	}

	// get all users to list
	jj, err := db.UserAll()
	if err != nil {
		return nil, err
	}

	var usrs []user.User
	for i := range jj {
		var u user.User
		err = json.Unmarshal(jj[i], &u)
		if err != nil {
			return nil, err
		}
		if u.Email != usr.Email {
			usrs = append(usrs, u)
		}
	}

	// make buffer to execute html into then pass buffer's bytes to Admin
	buf := &bytes.Buffer{}
	tmpl := template.Must(template.New("users").Parse(html + script))
	data := map[string]interface{}{
		"User":  usr,
		"Users": usrs,
	}

	err = tmpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return Admin(buf.Bytes())
}

var analyticsHTML = `
<div class="analytics">
<div class="card">
<div class="card-content">
    <p class="right">数据范围： {{ .from }} - {{ .to }} (UTC)</p>
    <div class="card-title">API 请求</div>
    <canvas id="analytics-chart"></canvas>
    <script>
    var target = document.getElementById("analytics-chart");
    Chart.defaults.global.defaultFontColor = '#212121';
    Chart.defaults.global.defaultFontFamily = "'Roboto', 'Helvetica Neue', 'Helvetica', 'Arial', 'sans-serif'";
    Chart.defaults.global.title.position = 'right';
    var chart = new Chart(target, {
        type: 'bar',
        data: {
            labels: [{{ range $date := .dates }} "{{ $date }}",  {{ end }}],
            datasets: [{
                type: 'line',
                label: '独立客户端',
                data: $.parseJSON({{ .unique }}),
                backgroundColor: 'rgba(76, 175, 80, 0.2)',
                borderColor: 'rgba(76, 175, 80, 1)',
                borderWidth: 1
            },
            {
                type: 'bar',
                label: '总请求',
                data: $.parseJSON({{ .total }}),
                backgroundColor: 'rgba(33, 150, 243, 0.2)',
                borderColor: 'rgba(33, 150, 243, 1)',
                borderWidth: 1
            }]
        },
        options: {
            scales: {
                yAxes: [{
                    ticks: {
                        beginAtZero:true
                    }
                }]
            }
        }
    });
    </script>
</div>
</div>
</div>
`

// Dashboard returns the admin view with analytics dashboard
func Dashboard() ([]byte, error) {
	buf := &bytes.Buffer{}
	data, err := analytics.ChartData()
	if err != nil {
		return nil, err
	}

	tmpl := template.Must(template.New("analytics").Parse(analyticsHTML))
	err = tmpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return Admin(buf.Bytes())
}

var err400HTML = []byte(`
<div class="error-page e400 col s6">
<div class="card">
<div class="card-content">
    <div class="card-title"><b>400</b> 错误： 错误的请求！</div>
    <blockquote>对不起，您的请求无法完成！</blockquote>
</div>
</div>
</div>
`)

// Error400 creates a subview for a 400 error page
func Error400() ([]byte, error) {
	return Admin(err400HTML)
}

var err404HTML = []byte(`
<div class="error-page e404 col s6">
<div class="card">
<div class="card-content">
    <div class="card-title"><b>404</b> 错误: 对象未找到</div>
    <blockquote>对不起，您要找的页面无法找到！</blockquote>
</div>
</div>
</div>
`)

// Error404 creates a subview for a 404 error page
func Error404() ([]byte, error) {
	return Admin(err404HTML)
}

var err405HTML = []byte(`
<div class="error-page e405 col s6">
<div class="card">
<div class="card-content">
    <div class="card-title"><b>405</b> 错误: 方法不允许</div>
    <blockquote>对不起，请求方法不允许</blockquote>
</div>
</div>
</div>
`)

// Error405 creates a subview for a 405 error page
func Error405() ([]byte, error) {
	return Admin(err405HTML)
}

var err500HTML = []byte(`
<div class="error-page e500 col s6">
<div class="card">
<div class="card-content">
    <div class="card-title"><b>500</b> 错误: 服务器内部错误</div>
    <blockquote>对不起，不可预知的错误发生了！</blockquote>
</div>
</div>
</div>
`)

// Error500 creates a subview for a 500 error page
func Error500() ([]byte, error) {
	return Admin(err500HTML)
}

var errMessageHTML = `
<div class="error-page eMsg col s6">
<div class="card">
<div class="card-content">
    <div class="card-title"><b>错误：&nbsp;</b>%s</div>
    <blockquote>%s</blockquote>
</div>
</div>
</div>
`

// ErrorMessage is a generic error message container, similar to Error500() and
// others in this package, ecxept it expects the caller to provide a title and
// message to describe to a view why the error is being shown
func ErrorMessage(title, message string) ([]byte, error) {
	eHTML := fmt.Sprintf(errMessageHTML, title, message)
	return Admin([]byte(eHTML))
}
