import { describe, it, expect, vi } from 'vitest'
import { scrollToHeading } from '../scroll'

describe('scrollToHeading', () => {
  it('does nothing if element not found', () => {
    const spy = vi.spyOn(window, 'scrollTo')
    scrollToHeading('nonexistent-id')
    expect(spy).not.toHaveBeenCalled()
    spy.mockRestore()
  })
})