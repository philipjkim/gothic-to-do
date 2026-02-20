# 2026-02-11: 프로젝트 초기화 세션

## 프로젝트 개요

- **목표**: 블로그 포스트 "GoTHIC 스택으로 JavaScript 없이 모던 웹앱 만들기"에 나오는 To-Do 앱을 실제로 구현
- **레포**: `github.com/philipjkim/gothic-to-do` (GitHub에서 클론)
- **작업 방식**: 사용자가 블로그 포스트를 보며 직접 코드를 타이핑 (학습 목적). 에러 발생 시 Claude에게 질문

## GoTHIC 스택 구성

| 기술 | 역할 |
|------|------|
| Go + Gin | 웹 서버, 라우팅 |
| Templ | Go 타입 안전 HTML 템플릿 엔진 |
| HTMX | HTML 속성만으로 AJAX 요청 및 DOM 교체 |
| Alpine.js | 클라이언트 전용 상태 관리 (토글, 모달 등) |
| TailwindCSS + DaisyUI | 유틸리티 퍼스트 CSS + UI 컴포넌트 |
| air | Go 핫 리로드 |

## 완료된 작업

### 1. 개발 도구 버전 업데이트

모든 도구를 최신 버전으로 업데이트 완료:

| 도구 | 이전 | 최신 | 관리 도구 | 업데이트 명령 |
|------|------|------|-----------|---------------|
| Go | 1.25.4 | **1.26.0** | asdf | `asdf install golang 1.26.0 && asdf set --home golang 1.26.0` |
| Node | 22.20.0 | **25.6.1** | **Volta** (asdf 아님) | `volta install node@latest` |
| npm | 11.6.2 | **11.9.0** | **Volta** | `volta install npm@latest` |
| templ | 0.3.960 | **0.3.977** | **go install** (brew 아님) | `go install github.com/a-h/templ/cmd/templ@latest` |
| air | 1.63.4 | **1.64.5** | **go install** (brew 아님) | `go install github.com/air-verse/air@latest` |

**주의사항**: Go asdf 설치 시 macOS 내장 `sha256sum`이 GNU 호환이 아니라서 체크섬 검증 실패함. 해결: `PATH="/opt/homebrew/opt/coreutils/libexec/gnubin:$PATH" asdf install golang <version>`

### 2. 파일 생성/수정

#### `.tool-versions` (신규)
```
golang 1.26.0
```
- Go 버전만 asdf로 관리하므로 여기에만 기록
- Node/npm은 Volta로 관리 → `package.json`의 `volta` 섹션에 핀 (아직 package.json 미생성)

#### `README.md` (수정)
- 블로그 포스트 기반으로 내용 보강
- GoTHIC 스택 소개, 프로젝트 구조, 설치/실행 방법, 주요 명령어
- HTMX vs Alpine.js 사용 기준을 mermaid flowchart로 표현

#### `.gitignore` (수정)
기존 Go 기본 항목에 아래 추가:
```
# air
tmp/
build-errors.log

# Node
node_modules/

# Tailwind build output
static/css/output.css

# templ generated files
*_templ.go
```

### 3. 커스텀 스킬 생성 (`~/.claude/skills/`)

#### `commit-kr` (한글 커밋)
- 경로: `~/.claude/skills/commit-kr/SKILL.md`
- Conventional Commits prefix + 한글 메시지
- Co-Authored-By 절대 불포함
- 커밋 후 push 여부 질문

#### `commit-en` (영문 커밋)
- 경로: `~/.claude/skills/commit-en/SKILL.md`
- Conventional Commits prefix + scope(선택) + 영문 메시지
- Co-Authored-By 절대 불포함
- 커밋 후 push 여부 질문

## 아직 하지 않은 작업

### CLAUDE.md 생성
- `/init`으로 생성 예정
- **코드 작업 완료 후에 만드는 게 낫다**는 결론 (프로젝트 구조, 빌드 명령어, 컨벤션이 확정된 후)

### 코드 작업 (사용자가 직접 진행)
블로그 포스트 순서대로:
1. `go mod init` + Go 의존성 설치
2. `npm init` + Node 의존성 설치 (TailwindCSS v4 주의 - 블로그는 v3 기준)
3. `tailwind.config.js` + `static/css/input.css`
4. `.air.toml`
5. `internal/model/todo.go`
6. `internal/store/memory.go`
7. Templ 템플릿들 (`templates/` 하위)
8. `internal/handler/todo.go`
9. `cmd/server/main.go`
10. `Makefile`

### Volta로 Node/npm 버전 핀
`package.json` 생성 후 `volta pin node@25` / `volta pin npm@11`

### git commit + push
- 지금까지 변경사항 커밋 필요 (.tool-versions, README.md, .gitignore)
- `/commit-kr` 또는 `/commit-en` 스킬 사용 예정

## 블로그 포스트에서 발견한 잠재적 이슈

1. **Templ 버전 차이**: 블로그 작성 시점과 현재 templ v0.3.977 사이에 API 변경 가능성. `templ.KV` 등 함수 동작 확인 필요

2. **패키지 구조 불일치** ✅ 해결: 블로그에서 `templates/components/` 하위 파일이 `package templates`로 선언됨. Go에서 디렉토리 = 패키지이므로 하위 디렉토리에 넣으면 별도 패키지가 됨. → **모든 `.templ` 파일을 `templates/`에 flat하게 두는 것으로 해결.** 블로그 포스트도 flat 구조로 수정함

3. **핸들러 import 누락**: `handler/todo.go`에서 `model` 패키지와 `templ` 패키지 import 선언이 블로그 코드에 빠져있음

4. **TailwindCSS v4 호환성**: `npm install -D tailwindcss @tailwindcss/cli daisyui` 시 v4가 설치될 수 있음. v4는 `tailwind.config.js` 대신 CSS-first 설정 방식이므로 블로그의 v3 설정과 충돌

5. **air 설정 파일명** ✅ 해결: 블로그에서 `.air.toml`(dot prefix)이지만 프로젝트에 `air.toml`(dot 없음)으로 생성함. air는 기본적으로 `.air.toml`을 찾으므로 설정을 못 읽고 루트 디렉토리를 `go build .` 하려다 "no Go files" 에러 발생. → **파일명을 `.air.toml`로 rename하여 해결**

6. **air 무한 리로드** ✅ 해결: air의 `include_ext`에 `go`가 포함되어 있어 `templ generate`로 생성된 `_templ.go` 파일 변경을 감지 → 다시 빌드 → `_templ.go` 재생성 → 무한 루프. → **`exclude_regex`에 `"_templ.go"` 추가하여 해결**

## 사용자 선호사항

- 코드는 직접 타이핑 (학습 목적)
- 에러/문제 발생 시 Claude에게 질문
- 적극적인 제안을 원함
- 사소한 것까지 질문해도 OK
