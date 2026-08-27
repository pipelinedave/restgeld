import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AboutModal from '../components/AboutModal.vue'

describe('AboutModal', () => {
  it('rendert nichts wenn visible=false', () => {
    const wrapper = mount(AboutModal, { props: { visible: false } })
    expect(wrapper.find('.modal-overlay').exists()).toBe(false)
  })

  it('rendert Philosophie und Metadaten wenn visible=true', () => {
    const wrapper = mount(AboutModal, { props: { visible: true } })
    expect(wrapper.find('.brand-title').text()).toBe('restgeld')
    expect(wrapper.text()).toContain('Die Philosophie')
    expect(wrapper.text()).toContain('Key Features')
    expect(wrapper.find('.meta-label').text()).toContain('Version')
    expect(wrapper.find('.commit-link').exists()).toBe(true)
    expect(wrapper.text()).not.toContain('v0.')
    expect(wrapper.find('.github-link').attributes('href')).toBe('https://github.com/pipelinedave/restgeld')
  })

  it('emittet close bei klick auf Schliessen', async () => {
    const wrapper = mount(AboutModal, { props: { visible: true } })
    await wrapper.find('.close-btn').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
