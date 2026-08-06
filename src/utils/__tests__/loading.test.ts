import { describe, it, expect, beforeEach } from 'vitest'
import { isGlobalLoading, startLoading, stopLoading } from '../loading'

describe('loading', () => {
  beforeEach(() => {
    while (isGlobalLoading.value) {
      stopLoading()
    }
  })

  it('starts loading', () => {
    startLoading()
    expect(isGlobalLoading.value).toBe(true)
  })

  it('stops loading', () => {
    startLoading()
    stopLoading()
    expect(isGlobalLoading.value).toBe(false)
  })

  it('handles nested loading calls', () => {
    startLoading()
    startLoading()
    stopLoading()
    expect(isGlobalLoading.value).toBe(true)

    stopLoading()
    expect(isGlobalLoading.value).toBe(false)
  })

  it('does not go below zero', () => {
    stopLoading()
    stopLoading()
    stopLoading()
    expect(isGlobalLoading.value).toBe(false)
  })
})