/*
 * Copyright (C) 2014-2015 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package webprogress

import (
	"errors"
	"sync"
)

var (
	ErrPackageInstallInProgress = errors.New("package installion in progress")
)

type Status struct {
	status map[string]*WebProgress
	l      sync.Mutex
}

func New() *Status {
	return &Status{status: make(map[string]*WebProgress)}
}

// Add add pkg to the list of progresses to track, it is idempotent
func (i *Status) Add(pkg string) (*WebProgress, error) {
	i.l.Lock()
	defer i.l.Unlock()

	if _, ok := i.status[pkg]; ok {
		return nil, ErrPackageInstallInProgress
	}

	i.status[pkg] = NewWebProgress()

	return i.status[pkg], nil
}

// Remove removes pkg to the list of progresses to track, it is a no op
// to remove multiple times.
func (i *Status) Remove(pkg string) {
	i.l.Lock()
	defer i.l.Unlock()

	delete(i.status, pkg)
}

// Get returns a *WebProgress corresponding to pkg or nil if not tracked.
func (i *Status) Get(pkg string) (*WebProgress, bool) {
	i.l.Lock()
	defer i.l.Unlock()

	w, ok := i.status[pkg]
	return w, ok
}
