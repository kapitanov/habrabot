package rss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeHTMLEmptyString(t *testing.T) {
	s, err := normalizeHTML("")
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, s)
}

func TestNormalizeHTMLPlainText(t *testing.T) {
	input := "plain old regular test"
	s, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, s, input)
}

func TestNormalizeHTMLBasicMarkup(t *testing.T) {
	input := "plain <b>bold</b> <strong>strong</strong> <i>italic</i> <i>italic and <b>bold</b></i>"
	expected := "plain <b>bold</b> <strong>strong</strong> <i>italic</i> <i>italic and <b>bold</b></i>"
	actual, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}

func TestNormalizeHTMLHyperlinkMarkup(t *testing.T) {
	input := "plain <a href=\"http://web.site\">url text</a> foo bar " +
		"<a href=\"http://web.site?utm_source=test\">skip me</a> zed"
	expected := "plain url text foo bar zed"
	actual, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}

func TestNormalizeHTMLImageMarkup(t *testing.T) {
	input := "plain <img src=\"http://image.url\"> text"
	expected := "plain text"
	actual, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}

func TestNormalizeHTMLLineBreak(t *testing.T) {
	input := "<br>plain<br>text"
	expected := "plain\ntext"
	actual, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}

func TestExtractImageURLEmptyText(t *testing.T) {
	input := ""

	actual, err := extractImageURL(input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, actual)
}

func TestExtractImageURLNoImgTag(t *testing.T) {
	input := "foo bar <b>foo</b><i>bar</i>"

	actual, err := extractImageURL(input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, actual)
}

func TestExtractImageURLOneImgTag(t *testing.T) {
	input := "foo bar <b>foo</b><i>bar</i> <img src=\"http:/test.image\">"
	expected := "http:/test.image"

	actual, err := extractImageURL(input)
	if err != nil {
		t.Fatal(err)
	}

	if assert.NotNil(t, actual) {
		assert.Equal(t, expected, *actual)
	}
}

func TestExtractImageURLFewImgTags(t *testing.T) {
	input := "foo bar <b>foo</b><i>bar</i> <img src=\"http:/test.image1\"> test <img src=\"http:/test.image2\">"
	expected := "http:/test.image1"

	actual, err := extractImageURL(input)
	if err != nil {
		t.Fatal(err)
	}

	if assert.NotNil(t, actual) {
		assert.Equal(t, expected, *actual)
	}
}

func TestHabrPreview1(t *testing.T) {
	input := "Привет, Хабр!<br> <br> Предлагаю вашему вниманию перевод статьи " +
		"\"<a href=\"http://blog.cleancoder.com/uncle-bob/2018/08/13/TooClean.html\">Too Clean?</a>\"" +
		" автора Robert C. Martin (Uncle Bob).<br> <br> " +
		"<img src=\"https://habrastorage.org/getpro/habr/post_images/80c/878/4e1/80c8784e1d022238be1b0cf707a7fed6.jpg\" alt=\"image\">" +
		"<br> <br> Я только что посмотрел выступление <a href=\"https://twitter.com/sarahmei\">Сары Мэй</a>:" +
		" <a href=\"https://www.youtube.com/watch?v=8_UoDmJi7U8\">Жизнеспособный код</a>." +
		" Это было очень хорошо. Я полностью согласен с основными моментами ее выступления. " +
		"С другой стороны, темой ее выступления было то, что я раньше должным образом не рассматривал.<br> " +
		"<a href=\"https://habr.com/ru/post/476076/?utm_source=habrahabr&amp;utm_medium=rss&amp;utm_campaign=476076#habracut\">" +
		"Читать дальше →</a>"
	expected := "Привет, Хабр!\n" +
		"Предлагаю вашему вниманию перевод статьи \"Too Clean?\" автора Robert C. Martin (Uncle Bob).\n" +
		"Я только что посмотрел выступление Сары Мэй: Жизнеспособный код. " +
		"Это было очень хорошо. Я полностью согласен с основными моментами ее выступления. " +
		"С другой стороны, темой ее выступления было то, что я раньше должным образом не рассматривал."

	actual, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}

func TestHabrPreview2(t *testing.T) {
	input := "<p>На текущий момент есть большое разнообразие обратных прокси серверов. Я перечислю только парочку из них.</p>" +
		"<p>Nginx</p>" +
		"<p>Envoy</p>" +
		"<p>HAProxy</p>" +
		"<p>Traefik</p>" +
		"<p>Также у каждого уважающего себя клауд провайдера есть свой прокси сервер.</p>" +
		"<p>AWS Elastic LoadBalancer</p>" +
		"<p>Google Cloud Load Balancer</p>" +
		"<p>DigitalOcean Load Balancer</p>" +
		"<p>Azure load balancer</p> " +
		"<a href=\"https://habr.com/ru/post/538936/?utm_source=habrahabr&amp;utm_medium=rss&amp;utm_campaign=538936#habracut\">" +
		"Читать далее</a>"
	expected := "На текущий момент есть большое разнообразие обратных прокси серверов. Я перечислю только парочку из них.\n" +
		"Nginx\n" +
		"Envoy\n" +
		"HAProxy\n" +
		"Traefik\n" +
		"Также у каждого уважающего себя клауд провайдера есть свой прокси сервер.\n" +
		"AWS Elastic LoadBalancer\n" +
		"Google Cloud Load Balancer\n" +
		"DigitalOcean Load Balancer\n" +
		"Azure load balancer"

	actual, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}

func TestHabrPreview3(t *testing.T) {
	input := "<p>А что, если я скажу вам, что линтеры для Go можно создавать вот таким декларативным способом?</p>" +
		"<br>\n" +
		"<pre><code class=\"go\">func alwaysTrue(m dsl.Matcher) {\n" +
		"    m.Match(`strings.Count($_, $_) &gt;= 0`).Report(`always evaluates to true`)\n" +
		"    m.Match(`bytes.Count($_, $_) &gt;= 0`).Report(`always evaluates to true`)\n" +
		"}\n" +
		"\n" +
		"func replaceAll() {\n" +
		"    m.Match(`strings.Replace($s, $d, $w, $n)`).\n" +
		"        Where(m[\"n\"].Value.Int() &lt;= 0).\n" +
		"        Suggest(`strings.ReplaceAll($s, $d, $w)`)\n" +
		"}</code></pre><br>\n" +
		"<p>Год назад я уже рассказывал об утилите <a href=\"https://github.com/quasilyte/go-ruleguard\" " +
		"rel=\"nofollow noopener noreferrer\">ruleguard</a>. " +
		"Сегодня хотелось бы поделиться тем, что нового появилось за это время.</p><br>\n" +
		"<p>Основные нововведения:</p><br>\n" +
		"<ul>\n" +
		"<li>Поддержка установки наборов правил через <a href=\"https://github.com/golang/go/wiki/Modules\" " +
		"rel=\"nofollow noopener noreferrer\">Go модули</a> (bundles)</li>\n" +
		"<li>Программируемые фильтры (компилируются в байт-код)</li>\n" +
		"<li>Добавлен режим отладки фильтров</li>\n" +
		"<li>Появился хороший обучающий материал: <a href=\"https://go-ruleguard.github.io/by-example/\" " +
		"rel=\"nofollow noopener noreferrer\">ruleguard by example</a></li>\n" +
		"<li>У проекта появились <a href=\"https://github.com/grafana/grafana/pull/28419\" " +
		"rel=\"nofollow noopener noreferrer\">реальные пользователи</a> и внешние <a href=\"https://github.com/dgryski/semgrep-go\" " +
		"rel=\"nofollow noopener noreferrer\">наборы правил</a></li>\n" +
		"<li><a href=\"https://go-ruleguard.github.io/play/\" " +
		"rel=\"nofollow noopener noreferrer\">Онлайн песочница</a>, позволяющая попробовать ruleguard прямо в браузере</li>\n" +
		"</ul><br>\n<img title=\"Автор иллюстрации: Татьяна Уфимцева @leased_line\" " +
		"src=\"https://habrastorage.org/webt/jb/iy/a0/jbiya0ab6njechtwp9dwtxhmw24.jpeg\"> " +
		"<a href=\"https://habr.com/ru/post/538930/?utm_source=habrahabr&amp;utm_medium=rss&amp;utm_campaign=538930#habracut\">" +
		"Читать дальше &rarr;</a>"

	expected := "А что, если я скажу вам, что линтеры для Go можно создавать вот таким декларативным способом?\n" +
		"<pre language=\"go\">func alwaysTrue(m dsl.Matcher) {\n" +
		"  m.Match(`strings.Count($_, $_) >= 0`).Report(`always evaluates to true`)\n" +
		"  m.Match(`bytes.Count($_, $_) >= 0`).Report(`always evaluates to true`)\n" +
		"}\n" +
		"func replaceAll() {\n" +
		"  m.Match(`strings.Replace($s, $d, $w, $n)`).\n" +
		"    Where(m[\"n\"].Value.Int() <= 0).\n" +
		"    Suggest(`strings.ReplaceAll($s, $d, $w)`)\n" +
		"}\n" +
		"</pre>\n" +
		"Год назад я уже рассказывал об утилите ruleguard. " +
		"Сегодня хотелось бы поделиться тем, что нового появилось за это время.\n" +
		"Основные нововведения:\n" +
		"- Поддержка установки наборов правил через Go модули (bundles)\n" +
		"- Программируемые фильтры (компилируются в байт-код)\n" +
		"- Добавлен режим отладки фильтров\n" +
		"- Появился хороший обучающий материал: ruleguard by example\n" +
		"- У проекта появились реальные пользователи и внешние наборы правил\n" +
		"- Онлайн песочница, позволяющая попробовать ruleguard прямо в браузере"

	actual, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}
