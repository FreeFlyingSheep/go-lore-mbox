# go-lore-mbox

Golang cli for custom style lore.kernel.org.

## Features

- Simple and easy to use.
- Parsing mbox to generate readable HTML.
- Supports custom style.

## Examples

### Flags

```text
-c string
    css file (default "assets/style.css")
-j string
    js file (default "assets/tools.js")
-u string
    https://lore.kernel.org/xxx/xxx
```

### Using custom style

The generated HTML file has the following format:

```html
<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <title>Title</title>
    <style>
        /* Copy from the CSS file */
    </style>
</header>

<body>
    <div class="thread">
    <a href="link">Subject</a> Author
    <ul>
        <li>
            <a href="link">Subject</a> Author
        </li>
    </ul>
    </div>

    <div class="content">
        <div id="id" class="message">
            <div class="subject">Subject</div>
            <div class="date">Date</div>

            <div id="button-0" class="button">
                <a href="javascript:fold('0')">[-] Collapse</a>
            </div>
            <div id="fold-0" class="message-header fold">
                <div id="button-1" class="button">
                    <a href="javascript:fold('1')">[-] Collapse</a>
                </div>
                <div id="fold-1" class="from fold">
                    From:
                    <ul>
                        <li>Name <a href="mailto:Email">&lt;Email&gt;</a></li>
                    </ul>
                </div>

                <div id="button-2" class="button">
                    <a href="javascript:fold('2')">[-] Collapse</a>
                </div>
                <div id="fold-2" class="to fold">
                    To:
                    <ul>
                        <li>Name <a href="mailto:Email">&lt;Email&gt;</a></li>
                    </ul>
                </div>
                
                <div id="button-3" class="button">
                    <a href="javascript:fold('3')">[-] Collapse</a>
                </div>
                <div id="fold-3" class="cc fold">
                    Cc:
                    <ul>
                        <li>Name <a href="mailto:Email">&lt;Email&gt;</a></li>
                    </ul>
                </div>
            </div>

            <div id="button-4" class="button">
                <a href="javascript:fold('4')">[-] Collapse</a>
            </div>
            <div id="fold-4" class="message-body fold">
                <div id="button-5" class="button">
                    <a href="javascript:fold('5')">[-] Collapse</a>
                </div>
                <div id="fold-5" class="text fold">
                    <pre>Code</pre>
                </div>
            </div>
        </div>
    </div>
    <script>
        /* Copy from the JS file */
    </script>
</body>

</html>
```

The `message-body` will contain the following CSS classes:

- `text`
- `start`
- `end`
- `before`
- `after`
- `change`
- `diff`
- `index`
- `add`
- `del`
- `quote`

You can specify custom CSS and JS files.
