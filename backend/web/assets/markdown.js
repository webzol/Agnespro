// Tiny markdown renderer for the embedded docs page.
// Supports: headings, paragraphs, lists, code blocks, inline code,
// blockquotes, tables, horizontal rules, links, bold, italic.
// Not a full CommonMark implementation; sized for our USER_GUIDE.md.

(function (root) {
  function escapeHTML(s) {
    return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
  }

  function renderInline(s) {
    s = escapeHTML(s);
    // links
    s = s.replace(/\[([^\]]+)\]\(([^)]+)\)/g, function (_m, text, url) {
      var safe = url.replace(/"/g, "&quot;");
      return '<a href="' + safe + '" target="_blank" rel="noopener">' + text + '</a>';
    });
    // bold
    s = s.replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");
    s = s.replace(/__([^_]+)__/g, "<strong>$1</strong>");
    // italic
    s = s.replace(/\*([^*]+)\*/g, "<em>$1</em>");
    s = s.replace(/_([^_]+)_/g, "<em>$1</em>");
    // inline code
    s = s.replace(/`([^`]+)`/g, "<code>$1</code>");
    return s;
  }

  function renderTable(lines) {
    // First two lines: header and separator (| --- | --- |)
    var header = lines[0].split("|").map(function (c) { return c.trim(); }).filter(Boolean);
    var alignRow = lines[1].split("|").map(function (c) { return c.trim(); }).filter(Boolean);
    var align = alignRow.map(function (c) {
      if (/^:-+:$/.test(c)) return "center";
      if (/^-+:$/.test(c)) return "right";
      if (/^:-+$/.test(c)) return "left";
      return "";
    });
    var body = lines.slice(2).map(function (l) {
      return l.split("|").map(function (c) { return c.trim(); }).filter(function (c) { return c !== ""; });
    });
    var out = "<table><thead><tr>";
    header.forEach(function (h, i) {
      var a = align[i] ? ' style="text-align:' + align[i] + '"' : "";
      out += "<th" + a + ">" + renderInline(h) + "</th>";
    });
    out += "</tr></thead><tbody>";
    body.forEach(function (row) {
      out += "<tr>";
      row.forEach(function (c, i) {
        var a = align[i] ? ' style="text-align:' + align[i] + '"' : "";
        out += "<td" + a + ">" + renderInline(c) + "</td>";
      });
      out += "</tr>";
    });
    out += "</tbody></table>";
    return out;
  }

  function render(md) {
    if (!md) return "";
    var lines = md.replace(/\r\n/g, "\n").split("\n");
    var html = [];
    var i = 0;

    function flushParagraph(buf) {
      if (buf.length === 0) return;
      html.push("<p>" + renderInline(buf.join(" ")) + "</p>");
      buf.length = 0;
    }

    var paraBuf = [];
    while (i < lines.length) {
      var line = lines[i];

      // Code fence
      if (/^```/.test(line)) {
        flushParagraph(paraBuf);
        var lang = line.replace(/^```/, "").trim();
        var codeLines = [];
        i++;
        while (i < lines.length && !/^```/.test(lines[i])) {
          codeLines.push(lines[i]);
          i++;
        }
        i++; // skip closing ```
        var cls = lang ? ' class="lang-' + lang + '"' : "";
        html.push("<pre" + cls + "><code>" + escapeHTML(codeLines.join("\n")) + "</code></pre>");
        continue;
      }

      // Heading
      var h = line.match(/^(#{1,6})\s+(.+)$/);
      if (h) {
        flushParagraph(paraBuf);
        var level = h[1].length;
        html.push("<h" + level + ">" + renderInline(h[2]) + "</h" + level + ">");
        i++;
        continue;
      }

      // Horizontal rule
      if (/^-{3,}\s*$/.test(line) || /^_{3,}\s*$/.test(line) || /^\*{3,}\s*$/.test(line)) {
        flushParagraph(paraBuf);
        html.push("<hr>");
        i++;
        continue;
      }

      // Table (must have | in line and next line be separator)
      if (line.indexOf("|") >= 0 && i + 1 < lines.length && /^\s*\|?\s*:?-+:?\s*(\|\s*:?-+:?\s*)+\|?\s*$/.test(lines[i + 1])) {
        flushParagraph(paraBuf);
        var tlines = [line, lines[i + 1]];
        i += 2;
        while (i < lines.length && lines[i].indexOf("|") >= 0 && lines[i].trim() !== "") {
          tlines.push(lines[i]);
          i++;
        }
        html.push(renderTable(tlines));
        continue;
      }

      // Blockquote
      if (/^>\s?/.test(line)) {
        flushParagraph(paraBuf);
        var quoteLines = [];
        while (i < lines.length && /^>\s?/.test(lines[i])) {
          quoteLines.push(lines[i].replace(/^>\s?/, ""));
          i++;
        }
        html.push("<blockquote>" + render(quoteLines.join("\n")) + "</blockquote>");
        continue;
      }

      // Unordered list
      if (/^[\-\*]\s+/.test(line)) {
        flushParagraph(paraBuf);
        var items = [];
        while (i < lines.length && /^[\-\*]\s+/.test(lines[i])) {
          items.push(lines[i].replace(/^[\-\*]\s+/, ""));
          i++;
        }
        html.push("<ul>" + items.map(function (it) { return "<li>" + renderInline(it) + "</li>"; }).join("") + "</ul>");
        continue;
      }

      // Ordered list
      if (/^\d+\.\s+/.test(line)) {
        flushParagraph(paraBuf);
        var oItems = [];
        while (i < lines.length && /^\d+\.\s+/.test(lines[i])) {
          oItems.push(lines[i].replace(/^\d+\.\s+/, ""));
          i++;
        }
        html.push("<ol>" + oItems.map(function (it) { return "<li>" + renderInline(it) + "</li>"; }).join("") + "</ol>");
        continue;
      }

      // Empty line -> paragraph break
      if (line.trim() === "") {
        flushParagraph(paraBuf);
        i++;
        continue;
      }

      paraBuf.push(line);
      i++;
    }
    flushParagraph(paraBuf);
    return html.join("\n");
  }

  root.MiniMark = { render: render };
})(window);