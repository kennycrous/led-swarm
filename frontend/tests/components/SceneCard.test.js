import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import SceneCard from '../../src/lib/components/SceneCard.svelte';

describe('SceneCard.svelte', () => {
  const mockScene = {
    id: 'scene-1',
    name: 'Cyberpunk Neon',
    icon: 'Sparkles',
    createdAt: '2026-08-27'
  };

  it('renders scene name and icon', () => {
    const { getByText } = render(SceneCard, {
      props: {
        scene: mockScene
      }
    });

    expect(getByText('Cyberpunk Neon')).toBeDefined();
  });

  it('triggers onApply when Apply Scene button is clicked', async () => {
    const handleApply = vi.fn();
    const { getByText } = render(SceneCard, {
      props: {
        scene: mockScene,
        onApply: handleApply
      }
    });

    const applyBtn = getByText('Apply Scene');
    await fireEvent.click(applyBtn);

    expect(handleApply).toHaveBeenCalledWith('scene-1');
  });
});
