# Snake 项目设计

## 目标

Snake 是一个基于 Bubble Tea v2 和 Lip Gloss v2 的终端贪吃蛇游戏。项目优先保证结构清晰、方便扩展，并让游戏核心逻辑不依赖具体 TUI 框架。

## 架构分层

```text
cmd/snake
  程序入口，只负责启动 TUI 程序。

internal/game
  游戏领域层，负责蛇、食物、地图、碰撞、得分、暂停和重开等纯规则。
  这一层不依赖 Bubble Tea 或 Lip Gloss，后续可以复用到 Web、测试、AI 控制器等场景。

internal/tui
  终端交互层，负责 Bubble Tea 的 Model、消息、按键、定时器和窗口尺寸适配。
  它调用 internal/game 驱动游戏状态变化。

internal/ui
  终端渲染层，负责 Lip Gloss 样式、棋盘绘制、彩色蛇身、状态栏和帮助文本。
  这一层只把 game.State 渲染成字符串，不修改游戏规则。
```

## 目录规划

```text
.
├── cmd/snake/main.go
├── docs/design.md
├── internal/game
│   ├── direction.go
│   ├── engine.go
│   ├── point.go
│   └── state.go
├── internal/tui
│   ├── app.go
│   ├── keys.go
│   └── messages.go
└── internal/ui
    ├── palette.go
    ├── renderer.go
    └── styles.go
```

## 核心模型

- `game.Point`：棋盘坐标。
- `game.Direction`：蛇的移动方向，负责反向移动保护。
- `game.State`：当前棋盘、蛇身、食物、分数、游戏状态和帧号。
- `game.Engine`：封装随机数和规则操作，提供 `Tick`、`Turn`、`TogglePause`、`Reset` 等方法。

## 交互设计

- 方向键或 `WASD` / `HJKL` 控制方向。
- `Space` 或 `P` 暂停和继续。
- `R` 重开。
- `Q`、`Esc` 或 `Ctrl+C` 退出。

## 视觉设计

- 棋盘使用 Lip Gloss 绘制边框和背景。
- 食物使用醒目的红粉色。
- 蛇头和蛇身分开渲染。
- 蛇身颜色由帧号和身体索引共同决定，tick 推进时颜色沿蛇身流动，形成游动变色效果。

## 扩展点

- 新增难度：调整 `game.Config` 的地图大小和 tick 间隔。
- 新增关卡：在 `game.Engine` 中加入障碍物集合。
- 新增 AI：实现外部控制器，调用 `game.Engine.Turn`。
- 新增渲染主题：扩展 `internal/ui` 的调色板和样式，不影响游戏规则。
