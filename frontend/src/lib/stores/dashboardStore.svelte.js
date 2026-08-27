// Svelte 5 Reactive Dashboard Store ($state runes)

let dashboardItems = $state([]);
let panels = $state([]);

const isWails = typeof window !== 'undefined' && window.runtime !== undefined;

export function getDashboardStore() {
  return {
    get dashboardItems() { return dashboardItems; },
    get panels() { return panels; },

    async init() {
      await loadDashboardItems();
      await loadDashboardPanels();
    },

    isPinned(itemId) {
      const found = dashboardItems.find(i => i.itemId === itemId);
      return found ? found.isPinned : true;
    },

    getSize(itemId) {
      const found = dashboardItems.find(i => i.itemId === itemId);
      return found?.size || 'normal';
    },

    getPanelId(itemId) {
      const found = dashboardItems.find(i => i.itemId === itemId);
      return found?.panelId || '';
    },

    async setSize(itemId, size) {
      // Optimistic local reactivity update
      const found = dashboardItems.find(i => i.itemId === itemId);
      if (found) {
        found.size = size;
      } else {
        dashboardItems.push({ itemId, itemType: 'device', size, isPinned: true, panelId: '' });
      }
      dashboardItems = [...dashboardItems];

      try {
        let updatedItem = null;
        if (isWails && window.go?.main?.App?.SetDashboardItemSize) {
          updatedItem = await window.go.main.App.SetDashboardItemSize(itemId, size);
        } else {
          const res = await fetch('/api/v1/dashboard/size', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ itemId, size })
          });
          if (res.ok) {
            updatedItem = await res.json();
          }
        }

        if (updatedItem) {
          upsertItem(updatedItem);
        }
      } catch (e) {
        console.error('Failed to set dashboard item size:', e);
      }
    },

    async setPanelId(itemId, panelId) {
      // Optimistic local update
      const found = dashboardItems.find(i => i.itemId === itemId);
      if (found) {
        found.panelId = panelId;
      } else {
        dashboardItems.push({ itemId, itemType: 'device', size: 'normal', isPinned: true, panelId });
      }
      dashboardItems = [...dashboardItems];

      try {
        let updatedItem = null;
        if (isWails && window.go?.main?.App?.SetDashboardItemPanel) {
          updatedItem = await window.go.main.App.SetDashboardItemPanel(itemId, panelId);
        } else {
          const res = await fetch('/api/v1/dashboard/panel', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ itemId, panelId })
          });
          if (res.ok) {
            updatedItem = await res.json();
          }
        }

        if (updatedItem) {
          upsertItem(updatedItem);
        }
      } catch (e) {
        console.error('Failed to set dashboard item panel:', e);
      }
    },

    async togglePin(itemId, itemType) {
      const currentPinned = this.isPinned(itemId);
      const targetPinned = !currentPinned;

      try {
        let updatedItem = null;
        if (isWails && window.go?.main?.App?.PinDashboardItem) {
          updatedItem = await window.go.main.App.PinDashboardItem(itemId, itemType, targetPinned);
        } else {
          const res = await fetch('/api/v1/dashboard/pin', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ itemId, itemType, isPinned: targetPinned })
          });
          if (res.ok) {
            updatedItem = await res.json();
          }
        }

        if (updatedItem) {
          upsertItem(updatedItem);
        }
      } catch (e) {
        console.error('Failed to toggle dashboard pin:', e);
      }
    },

    async addPanel(title) {
      if (!title.trim()) return;
      try {
        let newPanel = null;
        if (isWails && window.go?.main?.App?.AddDashboardPanel) {
          newPanel = await window.go.main.App.AddDashboardPanel(title.trim());
        } else {
          const res = await fetch('/api/v1/dashboard/panels', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ title: title.trim() })
          });
          if (res.ok) {
            newPanel = await res.json();
          }
        }

        if (newPanel) {
          panels = [...panels.filter(p => p.id !== newPanel.id), newPanel];
        } else {
          const fallback = { id: 'panel-' + Date.now(), title: title.trim() };
          panels = [...panels, fallback];
        }
      } catch (e) {
        console.error('Error adding panel:', e);
        const fallback = { id: 'panel-' + Date.now(), title: title.trim() };
        panels = [...panels, fallback];
      }
    },

    async deletePanel(panelId) {
      panels = panels.filter(p => p.id !== panelId);

      try {
        if (isWails && window.go?.main?.App?.DeleteDashboardPanel) {
          await window.go.main.App.DeleteDashboardPanel(panelId);
        } else {
          await fetch('/api/v1/dashboard/panels/' + panelId, {
            method: 'DELETE'
          });
        }
      } catch (e) {
        console.error('Error deleting panel:', e);
      }
    }
  };
}

async function loadDashboardItems() {
  try {
    let list = [];
    if (isWails && window.go?.main?.App?.GetDashboardItems) {
      list = await window.go.main.App.GetDashboardItems();
    } else {
      const res = await fetch('/api/v1/dashboard/items');
      if (res.ok) list = await res.json();
    }
    if (list) dashboardItems = list;
  } catch (e) {
    console.error('Error loading dashboard items:', e);
  }
}

async function loadDashboardPanels() {
  try {
    let list = [];
    if (isWails && window.go?.main?.App?.GetDashboardPanels) {
      list = await window.go.main.App.GetDashboardPanels();
    } else {
      const res = await fetch('/api/v1/dashboard/panels');
      if (res.ok) list = await res.json();
    }
    if (list && Array.isArray(list)) panels = list;
  } catch (e) {
    console.error('Error loading dashboard panels:', e);
  }
}

function upsertItem(newItem) {
  const idx = dashboardItems.findIndex(i => i.itemId === newItem.itemId);
  if (idx >= 0) {
    const current = dashboardItems[idx];
    dashboardItems[idx] = {
      ...current,
      ...newItem,
      panelId: (newItem.panelId !== undefined && newItem.panelId !== '') ? newItem.panelId : (current.panelId || ''),
      size: newItem.size || current.size || 'normal'
    };
  } else {
    dashboardItems.push(newItem);
  }
  dashboardItems = [...dashboardItems];
}
