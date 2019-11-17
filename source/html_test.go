package source

import "testing"

func TestNormalizeHTMLEmptyString(t *testing.T) {
	s, err := normalizeHTML("")
	if err != nil {
		t.Fatal(err)
	}

	if s != "" {
		t.Fatalf("expected empty string, got \"%s\"", s)
	}
}

func TestNormalizeHTMLPlainText(t *testing.T) {
	input := "plain old regular test"
	s, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	if s != input {
		t.Fatalf("expected \"%s\", got \"%s\"", input, s)
	}
}

func TestNormalizeHTMLBasicMarkup(t *testing.T) {
	input := "plain <b>bold</b> <strong>strong</strong> <i>italic</i> <i>italic and <b>bold</b></i>"
	expected := "plain bold strong italic italic and bold"
	s, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	if s != expected {
		t.Fatalf("expected \"%s\", got \"%s\"", expected, s)
	}
}

func TestNormalizeHTMLHyperlinkMarkup(t *testing.T) {
	input := "plain <a href=\"http://web.site\">url text</a> foo bar " +
		"<a href=\"http://web.site?utm_source=test\">skip me</a> zed"
	expected := "plain url text foo bar zed"
	s, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	if s != expected {
		t.Fatalf("expected \"%s\", got \"%s\"", expected, s)
	}
}

func TestNormalizeHTMLImageMarkup(t *testing.T) {
	input := "plain <img src=\"http://image.url\"> text"
	expected := "plain text"
	s, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	if s != expected {
		t.Fatalf("expected \"%s\", got \"%s\"", expected, s)
	}
}

func TestNormalizeHTMLLineBreak(t *testing.T) {
	input := "<br>plain<br>text"
	expected := "plain\ntext"
	s, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	if s != expected {
		t.Fatalf("expected \"%s\", got \"%s\"", expected, s)
	}
}

func TestExtractImageURLEmptyText(t *testing.T) {
	input := ""
	expected := ""

	actual, err := extractImageURL(input)
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("expected \"%s\", got \"%s\"", expected, actual)
	}
}

func TestExtractImageURLNoImgTag(t *testing.T) {
	input := "foo bar <b>foo</b><i>bar</i>"
	expected := ""

	actual, err := extractImageURL(input)
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("expected \"%s\", got \"%s\"", expected, actual)
	}
}

func TestExtractImageURLOneImgTag(t *testing.T) {
	input := "foo bar <b>foo</b><i>bar</i> <img src=\"http:/test.image\">"
	expected := "http:/test.image"

	actual, err := extractImageURL(input)
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("expected \"%s\", got \"%s\"", expected, actual)
	}
}

func TestExtractImageURLFewImgTags(t *testing.T) {
	input := "foo bar <b>foo</b><i>bar</i> <img src=\"http:/test.image1\"> test <img src=\"http:/test.image2\">"
	expected := "http:/test.image1"

	actual, err := extractImageURL(input)
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("expected \"%s\", got \"%s\"", expected, actual)
	}
}

func TestHabrPreview(t *testing.T) {
	input := "Привет, Хабр!<br> <br> Предлагаю вашему вниманию перевод статьи " +
		"\"<a href=\"http://blog.cleancoder.com/uncle-bob/2018/08/13/TooClean.html\">Too Clean?</a>\"" +
		" автора Robert C. Martin (Uncle Bob).<br> <br> " +
		"<img src=\"https://habrastorage.org/getpro/habr/post_images/80c/878/4e1/80c8784e1d022238be1b0cf707a7fed6.jpg\" alt=\"image\">" +
		"<br> <br> Я только что посмотрел выступление <a href=\"https://twitter.com/sarahmei\">Сары Мэй</a>:" +
		" <a href=\"https://www.youtube.com/watch?v=8_UoDmJi7U8\">Жизнеспособный код</a>." +
		" Это было очень хорошо. Я полностью согласен с основными моментами ее выступления. " +
		"С другой стороны, темой ее выступления было то, что я раньше должным образом не рассматривал.<br> " +
		"<a href=\"https://habr.com/ru/post/476076/?utm_source=habrahabr&amp;utm_medium=rss&amp;utm_campaign=476076#habracut\">Читать дальше →</a>"
	expected := "Привет, Хабр!\n" +
		"Предлагаю вашему вниманию перевод статьи \"Too Clean?\" автора Robert C. Martin (Uncle Bob).\n" +
		"Я только что посмотрел выступление Сары Мэй: Жизнеспособный код. " +
		"Это было очень хорошо. Я полностью согласен с основными моментами ее выступления. " +
		"С другой стороны, темой ее выступления было то, что я раньше должным образом не рассматривал."

	actual, err := normalizeHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("expected \"%s\", got \"%s\"", expected, actual)
	}
}
