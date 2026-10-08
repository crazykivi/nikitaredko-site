import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { scrollToHeading } from '../scroll'

describe('scrollToHeading', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('does nothing if element not found', () => {
    const spy = vi.spyOn(window, 'scrollTo')
    scrollToHeading('nonexistent-id')
    expect(spy).not.toHaveBeenCalled()
  })

  it('scrolls to element with offset', () => {
    const el = document.createElement('div')
    el.id = 'test-heading'
    document.body.appendChild(el)

    vi.spyOn(el, 'getBoundingClientRect').mockReturnValue({
      top: 500, bottom: 0, left: 0, right: 0, width: 0, height: 0, x: 0, y: 0, toJSON: () => { }
    })

    Object.defineProperty(window, 'scrollY', { value: 100, writable: true, configurable: true })
    const spy = vi.spyOn(window, 'scrollTo')

    scrollToHeading('test-heading', 50)

    expect(spy).toHaveBeenCalledWith({
      top: 550, // 500 (top) + 100 (scrollY) - 50 (offset)
      behavior: 'smooth'
    })
  })
})