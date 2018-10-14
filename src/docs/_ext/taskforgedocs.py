"""
Sphinx plugins for Taskforge documentation.
"""
import re
import logging

from docutils import nodes
from docutils.statemachine import ViewList
from sphinx.directives import CodeBlock

logger = logging.getLogger(__name__)
# RE for option descriptions without a '--' prefix
simple_option_desc_re = re.compile(
    r'([-_a-zA-Z0-9]+)(\s*.*?)(?=,\s+(?:/|-|--)|$)')


def setup(app):
    app.add_node(
        ConsoleNode,
        html=(visit_console_html, None),
        latex=(visit_console_dummy, depart_console_dummy),
        man=(visit_console_dummy, depart_console_dummy),
        text=(visit_console_dummy, depart_console_dummy),
        texinfo=(visit_console_dummy, depart_console_dummy),
    )
    app.add_directive('console', ConsoleDirective)
    app.connect('html-page-context', html_page_context_hook)
    return {'parallel_read_safe': True}


class ConsoleNode(nodes.literal_block):
    """
    Custom node to override the visit/depart event handlers at registration
    time. Wrap a literal_block object and defer to it.
    """
    tagname = 'ConsoleNode'

    def __init__(self, litblk_obj):
        self.wrapped = litblk_obj

    def __getattr__(self, attr):
        if attr == 'wrapped':
            return self.__dict__.wrapped
        return getattr(self.wrapped, attr)


def visit_console_dummy(self, node):
    """Defer to the corresponding parent's handler."""
    self.visit_literal_block(node)


def depart_console_dummy(self, node):
    """Defer to the corresponding parent's handler."""
    self.depart_literal_block(node)


def visit_console_html(self, node):
    """Generate HTML for the console directive."""
    if self.builder.name in ('taskforgehtml', 'json') and node['win_console_text']:
        # Put a mark on the document object signaling the fact the directive
        # has been used on it.
        self.document._console_directive_used_flag = True
        uid = node['uid']
        self.body.append('''\
<div class="console-block" id="console-block-%(id)s">
<input class="c-tab-unix" id="c-tab-%(id)s-unix" type="radio" name="console-%(id)s" checked>
<label for="c-tab-%(id)s-unix" title="Linux/macOS">&#xf17c/&#xf179</label>
<input class="c-tab-win" id="c-tab-%(id)s-win" type="radio" name="console-%(id)s">
<label for="c-tab-%(id)s-win" title="Windows">&#xf17a</label>
<section class="c-content-unix" id="c-content-%(id)s-unix">\n''' % {'id': uid})
        try:
            self.visit_literal_block(node)
        except nodes.SkipNode:
            pass
        self.body.append('</section>\n')

        self.body.append('<section class="c-content-win" id="c-content-%(id)s-win">\n' % {'id': uid})
        win_text = node['win_console_text']
        highlight_args = {'force': True}
        if 'linenos' in node:
            linenos = node['linenos']
        else:
            linenos = win_text.count('\n') >= self.highlightlinenothreshold - 1

        def warner(msg):
            self.builder.warn(msg, (self.builder.current_docname, node.line))

        highlighted = self.highlighter.highlight_block(
            win_text, 'doscon', warn=warner, linenos=linenos, **highlight_args
        )
        self.body.append(highlighted)
        self.body.append('</section>\n')
        self.body.append('</div>\n')
        raise nodes.SkipNode
    else:
        self.visit_literal_block(node)


class ConsoleDirective(CodeBlock):
    """
    A reStructuredText directive which renders a two-tab code block in which
    the second tab shows a Windows command line equivalent of the usual
    Unix-oriented examples.
    """
    required_arguments = 0
    # The 'doscon' Pygments formatter needs a prompt like this. '>' alone
    # won't do it because then it simply paints the whole command line as a
    # grey comment with no highlighting at all.
    WIN_PROMPT = r'...\> '

    def run(self):

        def args_to_win(cmdline):
            changed = False
            out = []
            for token in cmdline.split():
                if token[:2] == './':
                    token = token[2:]
                    changed = True
                elif token[:2] == '~/':
                    token = '%HOMEPATH%\\' + token[2:]
                    changed = True
                elif token == 'make':
                    token = 'make.bat'
                    changed = True
                if '://' not in token and 'git' not in cmdline:
                    out.append(token.replace('/', '\\'))
                    changed = True
                else:
                    out.append(token)
            if changed:
                return ' '.join(out)
            return cmdline

        def cmdline_to_win(line):
            if line.startswith('# '):
                return 'REM ' + args_to_win(line[2:])
            if line.startswith('$ # '):
                return 'REM ' + args_to_win(line[4:])
            if line.startswith('$ ./manage.py'):
                return 'manage.py ' + args_to_win(line[13:])
            if line.startswith('$ manage.py'):
                return 'manage.py ' + args_to_win(line[11:])
            if line.startswith('$ ./runtests.py'):
                return 'runtests.py ' + args_to_win(line[15:])
            if line.startswith('$ ./'):
                return args_to_win(line[4:])
            if line.startswith('$ python3'):
                return 'py ' + args_to_win(line[9:])
            if line.startswith('$ python'):
                return 'py ' + args_to_win(line[8:])
            if line.startswith('$ '):
                return args_to_win(line[2:])
            return None

        def code_block_to_win(content):
            bchanged = False
            lines = []
            for line in content:
                modline = cmdline_to_win(line)
                if modline is None:
                    lines.append(line)
                else:
                    lines.append(self.WIN_PROMPT + modline)
                    bchanged = True
            if bchanged:
                return ViewList(lines)
            return None

        env = self.state.document.settings.env
        self.arguments = ['console']
        lit_blk_obj = super().run()[0]

        # Only do work when the taskforgehtml HTML Sphinx builder is being used,
        # invoke the default behavior for the rest.
        if env.app.builder.name not in ('taskforgehtml', 'json'):
            return [lit_blk_obj]

        lit_blk_obj['uid'] = '%s' % env.new_serialno('console')
        # Only add the tabbed UI if there is actually a Windows-specific
        # version of the CLI example.
        win_content = code_block_to_win(self.content)
        if win_content is None:
            lit_blk_obj['win_console_text'] = None
        else:
            self.content = win_content
            lit_blk_obj['win_console_text'] = super().run()[0].rawsource

        # Replace the literal_node object returned by Sphinx's CodeBlock with
        # the ConsoleNode wrapper.
        return [ConsoleNode(lit_blk_obj)]


def html_page_context_hook(app, pagename, templatename, context, doctree):
    # Put a bool on the context used to render the template. It's used to
    # control inclusion of console-tabs.css and activation of the JavaScript.
    # This way it's include only from HTML files rendered from reST files where
    # the ConsoleDirective is used.
    context['include_console_assets'] = getattr(doctree, '_console_directive_used_flag', False)
