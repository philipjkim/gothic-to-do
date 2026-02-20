# 2026-02-11: 구현 세션

## 개요

블로그 포스트를 따라 GoTHIC To-Do 앱의 전체 코드를 작성하고, 발견된 이슈들을 수정하여 정상 동작까지 확인한 세션.

## 커밋 히스토리

| 커밋 | 설명 |
|------|------|
| `a5be2b6` | 초기 설정 (.gitignore, README, .tool-versions) |
| `d0d1611` | Go/Node 프로젝트 설정 (go.mod, package.json, air.toml, input.css) |
| `17c77aa` | README 영문화 + README_kr.md 분리 |
| `702c8ad` | Todo 모델 + 인메모리 스토어 |
| `3fc6afd` | Templ 템플릿 (layout, index, components) |
| `cec8ec8` | Gin 핸들러 + 템플릿 flat 구조 전환 + Makefile + 블로그 수정 |
| `c10913e` | 서버 엔트리포인트 + air 설정 수정 + 한글 주석 영문화 |
| `3719325` | bin/ 을 .gitignore 에 추가 |

## 완료된 작업

### 1. 프로젝트 설정 수정

- **go.mod 모듈명 오타 수정**: `philiphjkim` → `philipjkim`
- **go.mod Go 버전 수정**: `1.25.0` → `1.26.0`
- **TailwindCSS v4 설정**: 블로그는 v3 기준이었지만 v4 유지 결정. `tailwind.config.js` 삭제, `input.css`를 v4 문법(`@import "tailwindcss"`, `@plugin "daisyui"`)으로 전환
- **Volta 핀**: `package.json`에 `node 25.6.1`, `npm 11.9.0` 추가
- **.gitignore**: `.idea/` 주석 해제, `bin/` 추가

### 2. README 영문화

- 기존 `README.md` → `README_kr.md`로 rename
- 새 `README.md`를 영문으로 작성 (국제 개발자 대상)
- 영문 README 상단에 한국어 README 링크 포함

### 3. 코드 작성 (사용자가 직접 타이핑)

블로그 포스트 순서대로 코드 작성 완료:

| 파일 | 역할 |
|------|------|
| `internal/model/todo.go` | Todo 구조체 (ID, Title, Completed, CreatedAt) |
| `internal/store/memory.go` | 인메모리 CRUD + sync.RWMutex 동시성 처리 |
| `templates/layout.templ` | HTML 레이아웃 (HTMX, Alpine.js CDN) |
| `templates/index.templ` | 메인 페이지 (Layout + 컴포넌트 조합) |
| `templates/todo_form.templ` | 테마 토글 (Alpine.js) + 할일 추가 폼 (HTMX) |
| `templates/todo_list.templ` | 할일 목록 컨테이너 |
| `templates/todo_item.templ` | 할일 아이템 (인라인 편집 Alpine.js + CRUD HTMX) |
| `internal/handler/todo.go` | Gin HTTP 핸들러 (Index, Create, Toggle, Update, Delete) |
| `cmd/server/main.go` | 서버 엔트리포인트 (Gin 라우팅, 정적 파일 서빙) |
| `Makefile` | dev/build/templ/css 명령어 |

### 4. 코드 리뷰에서 발견/수정한 이슈

| 이슈 | 해결 |
|------|------|
| `store/memory.go` `GetAll()`에서 `RLock()` 후 `Unlock()` 호출 | `RUnlock()`으로 수정 (사용자가 직접 수정) |
| `todo_list.templ` import 경로 `github.com/yourname/todo-app` | `github.com/philipjkim/gothic-to-do`로 수정 (사용자가 직접 수정) |
| `TodoItem` div ID가 `todo-todo-000001`로 중복 | `fmt.Sprintf` 제거, `todo.ID` 직접 사용 |
| `templates/components/`가 `package templates` 선언 | Go 패키지 규칙 위반. flat 구조로 전환 (모든 templ 파일을 `templates/`에) |
| 블로그 포스트 핸들러에 `templ`, `model` import 누락 | 블로그 포스트 수정 |

### 5. air 설정 이슈 해결

| 이슈 | 원인 | 해결 |
|------|------|------|
| "no Go files" 에러 | 파일명이 `air.toml`이었지만 air는 `.air.toml`을 찾음 | `.air.toml`로 rename |
| 무한 리로드 루프 | `templ generate`가 `_templ.go` 생성 → air가 `.go` 변경 감지 → 다시 빌드 | `exclude_regex`에 `"_templ.go"` 추가 |
| `build.bin` deprecated 경고 | air v1.64+ 에서 `bin` deprecated | `entrypoint`로 변경 |

### 6. 개발 환경 이슈 해결

| 이슈 | 원인 | 해결 |
|------|------|------|
| `templ -v` 실행 시 "No version is set for command templ" | `.zshrc`에서 `GOPATH`가 쉘 시작 시점에 고정되어 Go 1.25.4 경로로 설정됨 | `.zshrc`에서 하드코딩된 `GOPATH` 제거, `$HOME/go/bin`을 PATH에 추가 |
| `air -v` 동일 에러 | asdf shim이 `~/go/bin`보다 PATH에서 앞에 위치 | 터미널에서 `source ~/.zshrc` 후 정상 동작 |

### 7. 한글 → 영문 전환

- 모든 `.templ` 파일의 한글 주석 및 UI 텍스트를 영문으로 변환
- `cmd/server/main.go`, `internal/handler/todo.go` 한글 주석 영문화
- `.air.toml`, `Makefile` 한글 주석 영문화

## 최종 상태

- `make dev`로 개발 서버 정상 구동 확인
- `localhost:8080`에서 모든 CRUD 기능 정상 동작 확인 (추가, 완료 토글, 인라인 편집, 삭제)
- 테마 토글 (cupcake/light/dark) 정상 동작 확인

## 아직 하지 않은 작업

- `CLAUDE.md` 생성 (`/init` 으로 생성 예정)
- 테스트 코드 작성
