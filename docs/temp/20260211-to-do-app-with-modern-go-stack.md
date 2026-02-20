+++
title = "GoTHIC 스택으로 JavaScript 없이 모던 웹앱 만들기"
date = "2026-02-11T09:31:05+09:00"
#dateFormat = "2006-01-02" # This value can be configured for per-post date formatting
author = "philipjkim"
authorTwitter = "" #do not include @
tags = ["golang", "programming", "dev"]
keywords = ["golang", "programming", "dev"]
description = "Go 백엔드 개발자가 프론트엔드 프레임워크 없이도 인터랙티브한 웹 애플리케이션을 만들 수 있다면 어떨까요? 이 글에서는 **GoTHIC 스택**으로 JavaScript 로직을 최소화한 To-Do 앱을 처음부터 끝까지 만들어봅니다. 핫 리로드까지 `air`로 세팅하면 개발 경험도 꽤 쾌적합니다."
showFullContent = false
readingTime = false
hideComments = false
+++

*이 글은 Claude Opus 4.5 을 이용해 초안이 작성되었으며, 이후 퇴고를 거쳤습니다.*

Go 백엔드 개발자가 프론트엔드 프레임워크 없이도 인터랙티브한 웹 애플리케이션을 만들 수 있다면 어떨까요? 이 글에서는 **GoTHIC 스택**으로 JavaScript 로직을 최소화한 To-Do 앱을 처음부터 끝까지 만들어봅니다. 핫 리로드까지 `air`로 세팅하면 개발 경험도 꽤 쾌적합니다.

---

## GoTHIC 스택이란

**GoTHIC**은 **Go**(+Gin) + **T**(empl) + **H**(TMX) + **I**nteractive **C**omponents(Alpine.js + TailwindCSS + DaisyUI)의 축약어입니다. GoTH(Go + Templ + HTMX) 스택을 확장하여 클라이언트 인터랙션과 스타일링까지 아우르는 풀스택 조합입니다.

| 기술 | GoTHIC에서의 역할 |
|------|------|
| **Go + Gin** | **Go** — 웹 프레임워크. 라우팅, 미들웨어 등 |
| **Templ** | **T** — Go 타입 안전 HTML 템플릿 엔진 |
| **HTMX** | **H** — HTML 속성만으로 AJAX 요청, DOM 교체 |
| **Alpine.js** | **IC** — 최소한의 클라이언트 상태 관리 (토글, 모달 등) |
| **TailwindCSS + DaisyUI** | **IC** — 유틸리티 퍼스트 CSS + UI 컴포넌트 |
| **air** | 개발 편의 — Go 핫 리로드 도구 |

React, Vue, Svelte 같은 SPA 프레임워크는 강력하지만, 모든 프로젝트에 필요한 건 아닙니다. 특히 Go 백엔드 개발자 입장에서 프론트엔드 빌드 파이프라인을 별도로 관리하는 건 부담이 될 수 있습니다.

GoTHIC의 핵심 아이디어는 **서버가 HTML을 렌더링하고, 브라우저는 그걸 잘 보여주는 것**입니다. HTMX가 서버와의 통신을 담당하고, Alpine.js가 순수 클라이언트 사이드 인터랙션(모달 열기/닫기, 드롭다운 토글 등)을 처리합니다. 둘의 역할이 명확히 나뉘기 때문에 JavaScript를 직접 작성할 일이 거의 없습니다.

## 프로젝트 초기 세팅

### 디렉토리 구조

```
todo-app/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   │   └── todo.go
│   ├── model/
│   │   └── todo.go
│   └── store/
│       └── memory.go
├── templates/
│   ├── layout.templ
│   ├── index.templ
│   ├── todo_form.templ
│   ├── todo_item.templ
│   └── todo_list.templ
├── static/
│   └── css/
│       └── output.css
├── .air.toml
├── tailwind.config.js
├── package.json
├── go.mod
└── go.sum
```

### Go 모듈 초기화 및 의존성 설치

```bash
mkdir todo-app && cd todo-app
go mod init github.com/yourname/todo-app

# Go 의존성
go get github.com/gin-gonic/gin
go get github.com/a-h/templ

# Templ CLI 설치
go install github.com/a-h/templ/cmd/templ@latest

# air 설치
go install github.com/air-verse/air@latest

# Node 의존성 (TailwindCSS + DaisyUI)
npm init -y
npm install -D tailwindcss @tailwindcss/cli daisyui
```

### TailwindCSS + DaisyUI 설정

`tailwind.config.js`:

```js
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["light", "dark", "cupcake"],
  },
}
```

`static/css/input.css`:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

TailwindCSS 빌드 명령:

```bash
npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css --watch
```

> 개발 시 `--watch` 플래그로 실행해두면 `.templ` 파일 변경 시 CSS가 자동으로 재빌드됩니다.

### air 설정

`.air.toml`:

```toml
root = "."
tmp_dir = "tmp"

[build]
  # templ generate 후 go build
  cmd = "templ generate && go build -o ./tmp/main ./cmd/server"
  bin = "./tmp/main"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "node_modules"]
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  include_ext = ["go", "templ", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_error = true

[log]
  time = false

[misc]
  clean_on_exit = true
```

이제 `air` 명령 하나로 `.go`와 `.templ` 파일 변경을 감지하여 자동으로 templ 코드 생성 → Go 빌드 → 서버 재시작이 이루어집니다.

## 모델과 인메모리 스토어

### 모델 정의

`internal/model/todo.go`:

```go
package model

import "time"

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}
```

### 인메모리 스토어

`internal/store/memory.go`:

```go
package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/yourname/todo-app/internal/model"
)

type TodoStore struct {
	mu     sync.RWMutex
	todos  map[string]*model.Todo
	nextID int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos: make(map[string]*model.Todo),
	}
}

func (s *TodoStore) Create(title string) *model.Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	todo := &model.Todo{
		ID:        fmt.Sprintf("todo-%d", s.nextID),
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	s.todos[todo.ID] = todo
	return todo
}

func (s *TodoStore) GetAll() []*model.Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]*model.Todo, 0, len(s.todos))
	for _, t := range s.todos {
		todos = append(todos, t)
	}
	return todos
}

func (s *TodoStore) Toggle(id string) (*model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return nil, fmt.Errorf("todo not found: %s", id)
	}
	todo.Completed = !todo.Completed
	return todo, nil
}

func (s *TodoStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return fmt.Errorf("todo not found: %s", id)
	}
	delete(s.todos, id)
	return nil
}

func (s *TodoStore) Update(id, title string) (*model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return nil, fmt.Errorf("todo not found: %s", id)
	}
	todo.Title = title
	return todo, nil
}
```

`sync.RWMutex`로 동시성을 처리합니다. 프로덕션에서는 당연히 실제 DB를 쓰겠지만, 이 예제에서는 스택 자체에 집중하기 위해 인메모리로 충분합니다.

## Templ 템플릿 작성

Templ은 Go의 타입 시스템을 활용하는 템플릿 엔진입니다. `.templ` 파일을 작성하면 `templ generate` 명령으로 Go 코드가 생성됩니다. 컴파일 타임에 타입 체크가 되기 때문에 런타임 에러를 크게 줄일 수 있습니다.

### 레이아웃

`templates/layout.templ`:

```templ
package templates

templ Layout(title string) {
	<!DOCTYPE html>
	<html lang="ko" data-theme="cupcake">
	<head>
		<meta charset="UTF-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
		<title>{ title }</title>
		<link href="/static/css/output.css" rel="stylesheet"/>
		<!-- HTMX -->
		<script src="https://unpkg.com/htmx.org@2.0.4"></script>
		<!-- Alpine.js -->
		<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
	</head>
	<body class="min-h-screen bg-base-200">
		<div class="container mx-auto max-w-2xl px-4 py-8">
			{ children... }
		</div>
	</body>
	</html>
}
```

DaisyUI의 `data-theme` 속성으로 테마를 간단하게 전환할 수 있습니다. CDN으로 HTMX와 Alpine.js를 로드했는데, 프로덕션에서는 로컬 번들링을 권장합니다.

### 메인 페이지

`templates/index.templ`:

```templ
package templates

import "github.com/yourname/todo-app/internal/model"

templ Index(todos []*model.Todo) {
	@Layout("To-Do App") {
		<div class="flex flex-col gap-6">
			<!-- 헤더: Alpine.js 테마 토글 -->
			@ThemeToggleHeader()

			<!-- 할일 입력 폼 -->
			@TodoForm()

			<!-- 할일 목록 -->
			<div id="todo-list">
				@TodoList(todos)
			</div>
		</div>
	}
}
```

### Alpine.js 테마 토글 (클라이언트 전용 인터랙션)

`templates/todo_form.templ`:

```templ
package templates

// ThemeToggleHeader - Alpine.js가 빛나는 순간입니다.
// 서버 왕복 없이 클라이언트에서 즉시 반응해야 하는 UI에 Alpine.js를 씁니다.
templ ThemeToggleHeader() {
	<div class="navbar bg-base-100 rounded-box shadow-md"
		x-data="{ 
			theme: localStorage.getItem('theme') || 'cupcake',
			themes: ['cupcake', 'light', 'dark'],
			setTheme(t) {
				this.theme = t;
				document.documentElement.setAttribute('data-theme', t);
				localStorage.setItem('theme', t);
			}
		}"
		x-init="setTheme(theme)"
	>
		<div class="flex-1">
			<span class="text-xl font-bold px-4">📝 To-Do App</span>
		</div>
		<div class="flex-none">
			<div class="dropdown dropdown-end">
				<div tabindex="0" role="button" class="btn btn-ghost btn-sm gap-1">
					🎨 테마
					<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
					</svg>
				</div>
				<ul tabindex="0" class="dropdown-content z-10 menu p-2 shadow-lg bg-base-100 rounded-box w-40">
					<template x-for="t in themes" :key="t">
						<li>
							<button
								x-text="t"
								@click="setTheme(t)"
								:class="{ 'active': theme === t }"
								class="capitalize"
							></button>
						</li>
					</template>
				</ul>
			</div>
		</div>
	</div>
}

// TodoForm - HTMX로 서버에 할일 추가 요청
templ TodoForm() {
	<div class="card bg-base-100 shadow-md">
		<div class="card-body p-4">
			<form
				hx-post="/todos"
				hx-target="#todo-list"
				hx-swap="innerHTML"
				hx-on::after-request="this.reset()"
				class="flex gap-2"
			>
				<input
					type="text"
					name="title"
					placeholder="할 일을 입력하세요..."
					required
					class="input input-bordered flex-1"
					autocomplete="off"
				/>
				<button type="submit" class="btn btn-primary">추가</button>
			</form>
		</div>
	</div>
}
```

여기서 HTMX와 Alpine.js의 역할 분리가 명확하게 드러납니다.

- **테마 토글**: 순수 클라이언트 로직이므로 Alpine.js가 담당합니다. 서버에 요청할 이유가 없습니다.
- **할일 추가 폼**: 서버에 데이터를 보내고 응답 HTML로 목록을 교체해야 하므로 HTMX가 담당합니다.

### 할일 목록 및 아이템

`templates/todo_list.templ`:

```templ
package templates

import "github.com/yourname/todo-app/internal/model"
import "fmt"

templ TodoList(todos []*model.Todo) {
	if len(todos) == 0 {
		<div class="card bg-base-100 shadow-md">
			<div class="card-body items-center text-center py-12">
				<p class="text-base-content/50 text-lg">아직 할 일이 없어요. 위에서 추가해보세요!</p>
			</div>
		</div>
	} else {
		<div class="card bg-base-100 shadow-md">
			<div class="card-body p-4">
				<div class="flex flex-col divide-y divide-base-200">
					for _, todo := range todos {
						@TodoItem(todo)
					}
				</div>
			</div>
		</div>
	}
}

// TodoItem - Alpine.js 인라인 편집 + HTMX 서버 동기화
templ TodoItem(todo *model.Todo) {
	<div
		id={ fmt.Sprintf("todo-%s", todo.ID) }
		class="flex items-center gap-3 py-3 group"
		x-data="{ editing: false }"
	>
		<!-- 완료 토글: HTMX -->
		<input
			type="checkbox"
			class="checkbox checkbox-primary"
			if todo.Completed {
				checked
			}
			hx-patch={ fmt.Sprintf("/todos/%s/toggle", todo.ID) }
			hx-target="#todo-list"
			hx-swap="innerHTML"
		/>

		<!-- 보기 모드 -->
		<span
			x-show="!editing"
			@dblclick="editing = true"
			class={ "flex-1 cursor-pointer select-none", templ.KV("line-through opacity-50", todo.Completed) }
		>
			{ todo.Title }
		</span>

		<!-- 편집 모드: Alpine.js 토글 + HTMX 저장 -->
		<form
			x-show="editing"
			x-cloak
			hx-put={ fmt.Sprintf("/todos/%s", todo.ID) }
			hx-target="#todo-list"
			hx-swap="innerHTML"
			@htmx:after-request.window="editing = false"
			class="flex-1 flex gap-2"
		>
			<input
				type="text"
				name="title"
				value={ todo.Title }
				class="input input-bordered input-sm flex-1"
				x-ref="editInput"
				@keydown.escape="editing = false"
			/>
			<button type="submit" class="btn btn-sm btn-success">저장</button>
			<button type="button" @click="editing = false" class="btn btn-sm btn-ghost">취소</button>
		</form>

		<!-- 삭제 버튼: HTMX -->
		<button
			x-show="!editing"
			hx-delete={ fmt.Sprintf("/todos/%s", todo.ID) }
			hx-target="#todo-list"
			hx-swap="innerHTML"
			hx-confirm="정말 삭제하시겠어요?"
			class="btn btn-ghost btn-sm opacity-0 group-hover:opacity-100 transition-opacity text-error"
		>
			✕
		</button>
	</div>
}
```

`TodoItem`이 GoTHIC 스택의 조합을 가장 잘 보여주는 컴포넌트입니다.

- **더블클릭으로 인라인 편집 전환**: `x-data="{ editing: false }"`와 `@dblclick="editing = true"`는 Alpine.js가 처리합니다. UI 상태 토글이므로 서버가 관여할 필요가 없습니다.
- **편집 내용 저장**: `hx-put`으로 서버에 보내고 응답 HTML로 목록을 교체합니다. 서버 통신은 HTMX의 영역입니다.
- **체크박스 토글, 삭제**: 마찬가지로 HTMX가 서버와 통신합니다.
- **호버 시 삭제 버튼 노출**: `group-hover:opacity-100`은 순수 CSS(Tailwind)입니다.

## Gin 핸들러

`internal/handler/todo.go`:

```go
package handler

import (
	"net/http"
	"sort"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/yourname/todo-app/internal/model"
	"github.com/yourname/todo-app/internal/store"
	"github.com/yourname/todo-app/templates"
)

type TodoHandler struct {
	store *store.TodoStore
}

func NewTodoHandler(store *store.TodoStore) *TodoHandler {
	return &TodoHandler{store: store}
}

// render는 templ 컴포넌트를 gin.Context에 렌더링하는 헬퍼입니다.
func render(c *gin.Context, component templ.Component) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "render error: %v", err)
	}
}

func (h *TodoHandler) Index(c *gin.Context) {
	todos := h.sortedTodos()
	render(c, templates.Index(todos))
}

func (h *TodoHandler) Create(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	h.store.Create(title)
	render(c, templates.TodoList(h.sortedTodos()))
}

func (h *TodoHandler) Toggle(c *gin.Context) {
	id := c.Param("id")
	if _, err := h.store.Toggle(id); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	render(c, templates.TodoList(h.sortedTodos()))
}

func (h *TodoHandler) Update(c *gin.Context) {
	id := c.Param("id")
	title := c.PostForm("title")
	if title == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	if _, err := h.store.Update(id, title); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	render(c, templates.TodoList(h.sortedTodos()))
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.store.Delete(id); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	render(c, templates.TodoList(h.sortedTodos()))
}

func (h *TodoHandler) sortedTodos() []*model.Todo {
	todos := h.store.GetAll()
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].CreatedAt.After(todos[j].CreatedAt)
	})
	return todos
}
```

핸들러 코드를 보면 패턴이 보입니다. **모든 변경 요청의 응답이 전체 목록의 HTML 파셜**입니다. HTMX가 `#todo-list`를 이 응답으로 교체하기 때문에 클라이언트에서 상태를 관리할 필요가 없습니다. 서버가 곧 진실의 원천(Single Source of Truth)이 되는 셈입니다.

> 참고로 `render` 헬퍼에서 import하는 `templ`은 `github.com/a-h/templ` 패키지입니다. `templ generate`로 생성된 코드와 함께 사용됩니다.

## 메인 서버

`cmd/server/main.go`:

```go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yourname/todo-app/internal/handler"
	"github.com/yourname/todo-app/internal/store"
)

func main() {
	r := gin.Default()

	// 정적 파일 서빙
	r.Static("/static", "./static")

	// 인메모리 스토어 & 핸들러
	todoStore := store.NewTodoStore()
	todoHandler := handler.NewTodoHandler(todoStore)

	// 라우트
	r.GET("/", todoHandler.Index)
	r.POST("/todos", todoHandler.Create)
	r.PATCH("/todos/:id/toggle", todoHandler.Toggle)
	r.PUT("/todos/:id", todoHandler.Update)
	r.DELETE("/todos/:id", todoHandler.Delete)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
```

## 개발 서버 실행

터미널을 두 개 열어서 각각 실행합니다.

```bash
# 터미널 1: TailwindCSS 빌드 (watch 모드)
npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css --watch

# 터미널 2: Go 서버 (air 핫 리로드)
air
```

브라우저에서 `http://localhost:8080`에 접속하면 To-Do 앱이 동작합니다. `.go` 파일이나 `.templ` 파일을 수정하면 air가 자동으로 재빌드하고, Tailwind 클래스를 변경하면 CSS가 자동으로 재빌드됩니다.

## HTMX vs Alpine.js: 언제 무엇을 쓸까

GoTHIC 스택에서 가장 중요한 판단은 "이 인터랙션을 HTMX로 할까, Alpine.js로 할까"입니다. 기준은 단순합니다.

**HTMX를 쓰는 경우**:
- 서버의 데이터가 변경되어야 할 때 (CRUD 작업)
- 응답으로 받은 HTML로 DOM을 교체할 때
- 폼 제출, 페이지네이션, 무한 스크롤 등

**Alpine.js를 쓰는 경우**:
- 순수 클라이언트 UI 상태 관리 (토글, 모달, 드롭다운)
- 서버 왕복 없이 즉시 반응해야 할 때
- 테마 전환, 폼 유효성 검사, 애니메이션 제어 등

**둘 다 쓰는 경우** (이 예제의 인라인 편집):
- Alpine.js가 편집 모드 전환 (UI 상태)을 관리하고
- HTMX가 수정된 내용을 서버로 전송 (데이터 변경)합니다

```
사용자 액션 → 서버 필요? ─ Yes → HTMX
                        └ No  → Alpine.js (또는 순수 CSS)
```

## 한 걸음 더: Makefile로 개발 편의성 높이기

매번 여러 터미널에서 명령어를 실행하는 건 번거롭습니다. `Makefile`로 정리해봅시다.

```makefile
.PHONY: dev build css templ

# 개발 서버 (air + tailwind watch)
dev:
	@echo "Starting development server..."
	@make -j2 air css-watch

air:
	air

css-watch:
	npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css --watch

# 빌드
build: templ css
	go build -o ./bin/server ./cmd/server

templ:
	templ generate

css:
	npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css --minify
```

이제 `make dev` 한 번이면 모든 워치 프로세스가 동시에 실행됩니다.

## 마무리

GoTHIC 스택의 본질은 **서버 사이드 렌더링으로 회귀하되, 현대적 DX를 포기하지 않는 것**입니다.

- **Templ**이 타입 안전한 HTML을 생성하고
- **HTMX**가 페이지 리로드 없는 인터랙션을 제공하며
- **Alpine.js**가 서버 왕복이 불필요한 클라이언트 인터랙션을 처리하고
- **TailwindCSS + DaisyUI**가 빠른 스타일링을 가능하게 하고
- **air**가 쾌적한 개발 경험을 보장합니다

JavaScript 번들러도 없고, 프론트엔드 빌드 파이프라인도 단순하며, Go 타입 시스템의 안전망 위에서 개발할 수 있습니다. 풀스택 SPA가 필요하지 않은 프로젝트라면 GoTHIC을 한 번 시도해볼 만합니다.

모든 소스 코드는 위에서 다룬 파일들을 그대로 구성하면 동작합니다. `go mod tidy` 잊지 마시고, Happy hacking!
