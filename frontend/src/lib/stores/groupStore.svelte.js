// Svelte 5 Reactive Group & Scene Store ($state runes)

let groups = $state([]);
let scenes = $state([]);
let activeGroupId = $state(null);

const isWails = typeof window !== 'undefined' && window.runtime !== undefined;

export function getGroupStore() {
  return {
    get groups() {
      return groups;
    },
    get scenes() {
      return scenes;
    },
    get activeGroupId() {
      return activeGroupId;
    },
    set activeGroupId(id) {
      activeGroupId = id;
    },

    async init() {
      await loadGroups();
      await loadScenes();
    },

    async createGroup(name, description, deviceIds) {
      try {
        let g = null;
        if (isWails && window.go?.main?.App?.SaveGroup) {
          g = await window.go.main.App.SaveGroup('', name, description, deviceIds);
        } else {
          const res = await fetch('/api/v1/groups', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, description, deviceIds })
          });
          if (res.ok) {
            g = await res.json();
          } else {
            const errText = await res.text();
            throw new Error(errText || 'Failed to create group');
          }
        }
        if (g) {
          upsertGroup(g);
        }
        return g;
      } catch (e) {
        console.error('Failed to create group:', e);
        throw e;
      }
    },

    async deleteGroup(id) {
      try {
        if (isWails && window.go?.main?.App?.DeleteGroup) {
          await window.go.main.App.DeleteGroup(id);
        } else {
          await fetch(`/api/v1/groups?id=${encodeURIComponent(id)}`, {
            method: 'DELETE'
          });
        }
        groups = groups.filter((g) => g.id !== id);
        if (activeGroupId === id) activeGroupId = null;
      } catch (e) {
        console.error('Failed to delete group:', e);
      }
    },

    async setGroupState(groupId, statePayload) {
      try {
        if (isWails && window.go?.main?.App?.SetGroupState) {
          await window.go.main.App.SetGroupState(groupId, JSON.stringify(statePayload));
        } else {
          await fetch('/api/v1/groups/state', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ groupId, state: statePayload })
          });
        }
      } catch (e) {
        console.error('Failed to set group state:', e);
      }
    },

    async captureScene(name, icon) {
      try {
        let s = null;
        if (isWails && window.go?.main?.App?.CaptureScene) {
          s = await window.go.main.App.CaptureScene(name, icon);
        } else {
          const res = await fetch('/api/v1/scenes/capture', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, icon })
          });
          if (res.ok) {
            s = await res.json();
          } else {
            const errText = await res.text();
            throw new Error(errText || 'Failed to capture scene');
          }
        }
        if (s) {
          upsertScene(s);
        }
        return s;
      } catch (e) {
        console.error('Failed to capture scene:', e);
        throw e;
      }
    },

    async applyScene(id) {
      try {
        if (isWails && window.go?.main?.App?.ApplyScene) {
          await window.go.main.App.ApplyScene(id);
        } else {
          await fetch('/api/v1/scenes/apply', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id })
          });
        }
      } catch (e) {
        console.error('Failed to apply scene:', e);
      }
    },

    async deleteScene(id) {
      try {
        if (isWails && window.go?.main?.App?.DeleteScene) {
          await window.go.main.App.DeleteScene(id);
        } else {
          await fetch(`/api/v1/scenes?id=${encodeURIComponent(id)}`, {
            method: 'DELETE'
          });
        }
        scenes = scenes.filter((s) => s.id !== id);
      } catch (e) {
        console.error('Failed to delete scene:', e);
      }
    }
  };
}

async function loadGroups() {
  try {
    let list = [];
    if (isWails && window.go?.main?.App?.GetGroups) {
      list = await window.go.main.App.GetGroups();
    } else {
      const res = await fetch('/api/v1/groups');
      if (res.ok) list = await res.json();
    }
    if (list) groups = list;
  } catch (e) {
    console.error('Error loading groups:', e);
  }
}

async function loadScenes() {
  try {
    let list = [];
    if (isWails && window.go?.main?.App?.GetScenes) {
      list = await window.go.main.App.GetScenes();
    } else {
      const res = await fetch('/api/v1/scenes');
      if (res.ok) list = await res.json();
    }
    if (list) scenes = list;
  } catch (e) {
    console.error('Error loading scenes:', e);
  }
}

function upsertGroup(newGroup) {
  const idx = groups.findIndex((g) => g.id === newGroup.id);
  if (idx >= 0) {
    groups[idx] = newGroup;
  } else {
    groups = [...groups, newGroup];
  }
}

function upsertScene(newScene) {
  const idx = scenes.findIndex((s) => s.id === newScene.id);
  if (idx >= 0) {
    scenes[idx] = newScene;
  } else {
    scenes = [...scenes, newScene];
  }
}
