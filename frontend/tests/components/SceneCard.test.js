import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import SceneCard from '../../src/lib/components/SceneCard.svelte';

describe('SceneCard.svelte', () => {
  const mockScene = {
    id: 'scene-1',
    name: 'Cyberpunk Neon',
    icon: 'Sparkles',
    configJson: '{}'
  };

  it('renders scene name and icon', () => {
    const { getByText } = render(SceneCard, {
      scene: mockScene
    });

    expect(getByText('Cyberpunk Neon')).toBeDefined();
    expect(getByText('Apply Scene Snapshot')).toBeDefined();
  });

  it('triggers onApply when Apply Scene button is clicked', async () => {
    const onApply = vi.fn();
    const { getByText } = render(SceneCard, {
      scene: mockScene,
      onApply
    });

    const applyBtn = getByText('Apply Scene Snapshot');
    await fireEvent.click(applyBtn);

    expect(onApply).toHaveBeenCalledWith('scene-1');
  });
});
