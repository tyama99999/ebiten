// Copyright 2016 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build android ios

package mobile

import (
	"github.com/hajimehoshi/ebiten/internal/driver"
)

type pos struct {
	X int
	Y int
}

type Input struct {
	cursorX int
	cursorY int
	keys    map[driver.Key]struct{}
	runes   []rune
	touches map[int]pos
	ui      *UserInterface
}

func (i *Input) CursorPosition() (x, y int) {
	i.ui.m.RLock()
	defer i.ui.m.RUnlock()
	return i.ui.adjustPosition(i.cursorX, i.cursorY)
}

func (i *Input) GamepadIDs() []int {
	return nil
}

func (i *Input) GamepadSDLID(id int) string {
	return ""
}

func (i *Input) GamepadName(id int) string {
	return ""
}

func (i *Input) GamepadAxisNum(id int) int {
	return 0
}

func (i *Input) GamepadAxis(id int, axis int) float64 {
	return 0
}

func (i *Input) GamepadButtonNum(id int) int {
	return 0
}

func (i *Input) IsGamepadButtonPressed(id int, button driver.GamepadButton) bool {
	return false
}

func (i *Input) TouchIDs() []int {
	i.ui.m.RLock()
	defer i.ui.m.RUnlock()

	if len(i.touches) == 0 {
		return nil
	}

	var ids []int
	for id := range i.touches {
		ids = append(ids, id)
	}
	return ids
}

func (i *Input) TouchPosition(id int) (x, y int) {
	i.ui.m.RLock()
	defer i.ui.m.RUnlock()

	for tid, pos := range i.touches {
		if id == tid {
			return i.ui.adjustPosition(pos.X, pos.Y)
		}
	}
	return 0, 0
}

func (i *Input) RuneBuffer() []rune {
	return i.runes
}

func (i *Input) IsKeyPressed(key driver.Key) bool {
	i.ui.m.RLock()
	defer i.ui.m.RUnlock()

	if i.keys == nil {
		return false
	}
	_, ok := i.keys[key]
	return ok
}

func (i *Input) Wheel() (xoff, yoff float64) {
	return 0, 0
}

func (i *Input) IsMouseButtonPressed(key driver.MouseButton) bool {
	return false
}

func (i *Input) update(keys map[driver.Key]struct{}, runes []rune, touches []*Touch) {
	i.ui.m.Lock()
	defer i.ui.m.Unlock()

	i.keys = map[driver.Key]struct{}{}
	for k := range keys {
		i.keys[k] = struct{}{}
	}

	i.runes = make([]rune, len(runes))
	copy(i.runes, runes)

	i.touches = map[int]pos{}
	for _, t := range touches {
		i.touches[t.ID] = pos{
			X: t.X,
			Y: t.Y,
		}
	}
}

func (i *Input) ResetForFrame() {
	i.runes = nil
}
