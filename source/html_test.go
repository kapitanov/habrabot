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
	input := "plain <a href=\"http://web.site\">url text</a> foo bar"
	expected := "plain  foo bar"
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
	expected := "plain  text"
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
