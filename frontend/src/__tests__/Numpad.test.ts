import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Numpad from '../components/Numpad.vue'

describe('Numpad', () => {
  it('rendert nichts wenn visible=false', () => {
    const wrapper = mount(Numpad, { props: { visible: false } })
    expect(wrapper.find('.numpad-overlay').exists()).toBe(false)
  })

  it('rendert ziffern und steuerbuttons wenn visible=true', () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    expect(wrapper.text()).toContain('0')
    expect(wrapper.text()).toContain('1')
    expect(wrapper.text()).toContain('7')
    expect(wrapper.text()).toContain(',')
    expect(wrapper.text()).toContain('⌫')
    expect(wrapper.text()).toContain('Bestätigen')
    expect(wrapper.text()).toContain('Abbrechen')
  })

  it('zeigt display 0 initial', () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    expect(wrapper.find('.numpad-value').text()).toBe('0')
  })

  it('fuegt ziffer bei klick hinzu', async () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    await wrapper.findAll('.numpad-btn')[0].trigger('click') // 7
    expect(wrapper.find('.numpad-value').text()).toBe('7')
  })

  it('fuegt komma hinzu', async () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    await wrapper.findAll('.numpad-btn')[0].trigger('click') // 7
    await wrapper.findAll('.numpad-btn')[9].trigger('click') // ,
    expect(wrapper.find('.numpad-value').text()).toBe('7,')
    await wrapper.findAll('.numpad-btn')[9].trigger('click') // zweites , ignoriert
    expect(wrapper.find('.numpad-value').text()).toBe('7,')
  })

  it('loescht letztes zeichen', async () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    await wrapper.findAll('.numpad-btn')[0].trigger('click') // 7
    await wrapper.findAll('.numpad-btn')[1].trigger('click') // 8
    expect(wrapper.find('.numpad-value').text()).toBe('78')
    await wrapper.findAll('.numpad-btn')[11].trigger('click') // backspace
    expect(wrapper.find('.numpad-value').text()).toBe('7')
  })

  it('emittet cancel bei abbrechen', async () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    await wrapper.findAll('.numpad-btn')[12].trigger('click') // abbrechen
    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('bestätigen zeigt notiz-eingabe statt confirm-emit', async () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    await wrapper.findAll('.numpad-btn')[0].trigger('click') // 7
    await wrapper.findAll('.numpad-btn')[13].trigger('click') // bestätigen
    expect(wrapper.find('.note-input').exists()).toBe(true)
    expect(wrapper.find('.note-section').exists()).toBe(true)
    expect(wrapper.emitted('confirm')).toBeFalsy()
  })

  it('notiz eingeben und speichern emittet confirm mit amount + note', async () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    await wrapper.findAll('.numpad-btn')[0].trigger('click') // 7
    await wrapper.findAll('.numpad-btn')[4].trigger('click') // 5
    await wrapper.findAll('.numpad-btn')[9].trigger('click') // ,
    await wrapper.findAll('.numpad-btn')[10].trigger('click') // 0
    await wrapper.findAll('.numpad-btn')[13].trigger('click') // bestätigen

    const noteInput = wrapper.find('.note-input')
    await noteInput.setValue('Kaffee')
    await wrapper.findAll('.note-actions button')[1].trigger('click') // speichern

    expect(wrapper.emitted('confirm')).toBeTruthy()
    expect(wrapper.emitted('confirm')![0]).toEqual([75, 'Kaffee']) // 75,0 wird zu 75
  })

  it('zurueck vom notiz-modus zeigt numpad wieder', async () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    await wrapper.findAll('.numpad-btn')[0].trigger('click') // 7
    await wrapper.findAll('.numpad-btn')[13].trigger('click') // bestätigen
    expect(wrapper.find('.note-section').exists()).toBe(true)

    await wrapper.findAll('.note-actions button')[0].trigger('click') // zurück
    expect(wrapper.find('.numpad-grid').exists()).toBe(true)
    expect(wrapper.find('.note-section').exists()).toBe(false)
  })

  it('bestätigt nicht ohne eingabe', async () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    await wrapper.findAll('.numpad-btn')[13].trigger('click') // bestätigen ohne input
    expect(wrapper.emitted('confirm')).toBeFalsy()
    expect(wrapper.find('.numpad-grid').exists()).toBe(true) // immer noch im numpad
  })
})
