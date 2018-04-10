// Package config provides a content type to manage the Ponzu system's configuration
// settings for things such as its name, domain, HTTP(s) port, email, server defaults
// and backups.
package config

import (
	"github.com/ponzu-cms/ponzu/management/editor"
	"github.com/ponzu-cms/ponzu/system/item"
)

// Config represents the confirgurable options of the system
type Config struct {
	item.Item

	Name                    string   `json:"name"`
	Domain                  string   `json:"domain"`
	BindAddress             string   `json:"bind_addr"`
	HTTPPort                string   `json:"http_port"`
	HTTPSPort               string   `json:"https_port"`
	AdminEmail              string   `json:"admin_email"`
	ClientSecret            string   `json:"client_secret"`
	Etag                    string   `json:"etag"`
	DisableCORS             bool     `json:"cors_disabled"`
	DisableGZIP             bool     `json:"gzip_disabled"`
	DisableHTTPCache        bool     `json:"cache_disabled"`
	CacheMaxAge             int64    `json:"cache_max_age"`
	CacheInvalidate         []string `json:"cache"`
	BackupBasicAuthUser     string   `json:"backup_basic_auth_user"`
	BackupBasicAuthPassword string   `json:"backup_basic_auth_password"`
}

const (
	dbBackupInfo = `
		<p class="flow-text">���ݿⱸ����֤��</p>
		<p>���һ���û����������������������ݿⱸ���ļ���HTTP���ء�</p>
	`
)

// String partially implements item.Identifiable and overrides Item's String()
func (c *Config) String() string { return c.Name }

// MarshalEditor writes a buffer of html to edit a Post and partially implements editor.Editable
func (c *Config) MarshalEditor() ([]byte, error) {
	view, err := editor.Form(c,
		editor.Field{
			View: editor.Input("Name", c, map[string]string{
				"label":       "վ�����ƣ������ڲ�ʹ�ã�",
				"placeholder": "����һ��վ�����ƣ������ڲ�ʹ�ã�",
			}),
		},
		editor.Field{
			View: editor.Input("Domain", c, map[string]string{
				"label":       "���� ������SSL֤��ʱ��Ҫ��",
				"placeholder": "���磺 www.example.com �� example.com",
			}),
		},
		editor.Field{
			View: editor.Input("BindAddress", c, map[string]string{
				"type": "hidden",
			}),
		},
		editor.Field{
			View: editor.Input("HTTPPort", c, map[string]string{
				"type": "hidden",
			}),
		},
		editor.Field{
			View: editor.Input("HTTPSPort", c, map[string]string{
				"type": "hidden",
			}),
		},
		editor.Field{
			View: editor.Input("AdminEmail", c, map[string]string{
				"label": "����Ա�ʼ���ַ����Ϣ���Ѻ��ڲ���Ϣ�����ã�",
			}),
		},
		editor.Field{
			View: editor.Input("ClientSecret", c, map[string]string{
				"label":    "�ͻ������루������֤����һ���˱��������",
				"disabled": "true",
			}),
		},
		editor.Field{
			View: editor.Input("ClientSecret", c, map[string]string{
				"type": "hidden",
			}),
		},
		editor.Field{
			View: editor.Input("Etag", c, map[string]string{
				"label":    "Etagͷ��Ϣ�����ڻ�����ƣ�",
				"disabled": "true",
			}),
		},
		editor.Field{
			View: editor.Input("Etag", c, map[string]string{
				"type": "hidden",
			}),
		},
		editor.Field{
			View: editor.Checkbox("DisableCORS", c, map[string]string{
				"label": "���� CORS ��������" + c.Domain + "���Է��ʵ��������ݣ�",
			}, map[string]string{
				"true": "���� CORS",
			}),
		},
		editor.Field{
			View: editor.Checkbox("DisableGZIP", c, map[string]string{
				"label": "���� GZIP ������ѹ��GZIP�������������ٶȣ��������ĸ������",
			}, map[string]string{
				"true": "���� GZIP",
			}),
		},
		editor.Field{
			View: editor.Checkbox("DisableHTTPCache", c, map[string]string{
				"label": "���� HTTP ���� ����д 'Cache-Control' ͷ��Ϣ��",
			}, map[string]string{
				"true": "���� HTTP ����",
			}),
		},
		editor.Field{
			View: editor.Input("CacheMaxAge", c, map[string]string{
				"label": "HTTP�����Max-Ageֵ����λ���룬0 �� 259200��",
				"type":  "text",
			}),
		},
		editor.Field{
			View: editor.Checkbox("CacheInvalidate", c, map[string]string{
				"label": "�����ͬʱʹ����ʧЧ",
			}, map[string]string{
				"invalidate": "ʹ����ʧЧ",
			}),
		},
		editor.Field{
			View: []byte(dbBackupInfo),
		},
		editor.Field{
			View: editor.Input("BackupBasicAuthUser", c, map[string]string{
				"label":       "HTTP��֤�û���",
				"placeholder": "����һ���û���",
				"type":        "text",
			}),
		},
		editor.Field{
			View: editor.Input("BackupBasicAuthPassword", c, map[string]string{
				"label":       "HTTP��֤����",
				"placeholder": "����һ������",
				"type":        "password",
			}),
		},
	)
	if err != nil {
		return nil, err
	}

	open := []byte(`
	<div class="card">
		<div class="card-content">
			<div class="card-title">ϵͳ����</div>
		</div>
		<form action="/admin/configure" method="post">
	`)
	close := []byte(`</form></div>`)
	script := []byte(`
	<script>
		$(function() {
			// hide default fields & labels unnecessary for the config
			var fields = $('.default-fields');
			fields.css('position', 'relative');
			fields.find('input:not([type=submit])').remove();
			fields.find('label').remove();
			fields.find('button').css({
				position: 'absolute',
				top: '-10px',
				right: '0px'
			});

			var contentOnly = $('.content-only.__ponzu');
			contentOnly.hide();
			contentOnly.find('input, textarea, select').attr('name', '');

			// adjust layout of td so save button is in same location as usual
			fields.find('td').css('float', 'right');

			// stop some fixed config settings from being modified
			fields.find('input[name=client_secret]').attr('name', '');
		});
	</script>
	`)

	view = append(open, view...)
	view = append(view, close...)
	view = append(view, script...)

	return view, nil
}
