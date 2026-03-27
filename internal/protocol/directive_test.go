package protocol

import "testing"

func TestParseDirectiveSimple(t *testing.T) {
	dir, id, args := ParseDirective("::ok Tests passing")
	if dir != "ok" { t.Errorf("directive: %q", dir) }
	if id != "" { t.Errorf("id: %q", id) }
	if args != "Tests passing" { t.Errorf("args: %q", args) }
}

func TestParseDirectiveWithID(t *testing.T) {
	dir, id, args := ParseDirective("::ok id=build Build complete")
	if dir != "ok" { t.Errorf("directive: %q", dir) }
	if id != "build" { t.Errorf("id: %q", id) }
	if args != "Build complete" { t.Errorf("args: %q", args) }
}

func TestParseDirectiveBarWithID(t *testing.T) {
	dir, id, args := ParseDirective("::bar id=prog 75 Coverage")
	if dir != "bar" { t.Errorf("directive: %q", dir) }
	if id != "prog" { t.Errorf("id: %q", id) }
	if args != "75 Coverage" { t.Errorf("args: %q", args) }
}

func TestParseDirectiveStatus(t *testing.T) {
	dir, id, args := ParseDirective("::status id=db running Migrations")
	if dir != "status" { t.Errorf("directive: %q", dir) }
	if id != "db" { t.Errorf("id: %q", id) }
	if args != "running Migrations" { t.Errorf("args: %q", args) }
}

func TestParseDirectivePlainText(t *testing.T) {
	dir, id, args := ParseDirective("Just regular text")
	if dir != "" || id != "" || args != "" { t.Errorf("plain: dir=%q id=%q args=%q", dir, id, args) }
}

func TestParseDirectiveEmpty(t *testing.T) {
	dir, _, _ := ParseDirective("::")
	if dir != "" { t.Errorf("empty: %q", dir) }
}

func TestParseDirectiveNoArgs(t *testing.T) {
	dir, id, args := ParseDirective("::divider")
	if dir != "divider" { t.Errorf("dir: %q", dir) }
	if id != "" { t.Errorf("id: %q", id) }
	if args != "" { t.Errorf("args: %q", args) }
}

func TestParseDirectiveIDOnly(t *testing.T) {
	dir, id, args := ParseDirective("::spin id=deploy")
	if dir != "spin" || id != "deploy" || args != "" { t.Errorf("got dir=%q id=%q args=%q", dir, id, args) }
}

func TestParseDirectiveCaseInsensitive(t *testing.T) {
	dir, _, _ := ParseDirective("::OK Tests")
	if dir != "ok" { t.Errorf("case: %q", dir) }
}

func TestParseDirectiveKVWithID(t *testing.T) {
	dir, id, args := ParseDirective("::kv id=meta CPU=58% | Mem=76%")
	if dir != "kv" || id != "meta" || args != "CPU=58% | Mem=76%" { t.Errorf("got dir=%q id=%q args=%q", dir, id, args) }
}
