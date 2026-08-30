import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Numpad from '../components/Numpad.vue'

describe('Numpad / Ausgabe-Modal', () => {
  it('rendert nichts wenn visible=false', () => {
    const wrapper = mount(Numpad, { props: { visible: false } })
    expect(wrapper.find('.modal-overlay').exists()).toBe(false)
  })

  it('rendert formularelemente wenn visible=true', () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    expect(wrapper.find('.modal-overlay').exists()).toBe(true)
    expect(wrapper.find('h2').text()).toBe('Ausgabe buchen')

    const amountInput = wrapper.find<HTMLInputElement>('#expense-amount-input')
    expect(amountInput.exists()).toBe(true)
    expect(amountInput.attributes('inputmode')).toBe('decimal')

    const noteInput = wrapper.find<HTMLInputElement>('#expense-note-input')
    expect(noteInput.exists()).toBe(true)
    expect(noteInput.attributes('inputmode')).toBe('text')

    expect(wrapper.find('.btn-confirm').exists()).toBe(true)
    expect(wrapper.find('.btn-cancel').exists()).toBe(true)
  })

  it('deaktiviert speichern-button wenn betrag ungueltig oder leer ist', async () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    const submitBtn = wrapper.find<HTMLButtonElement>('.btn-confirm')
    expect(submitBtn.element.disabled).toBe(true)

    const amountInput = wrapper.find<HTMLInputElement>('#expense-amount-input')
    await amountInput.setValue('0')
    expect(submitBtn.element.disabled).toBe(true)

    await amountInput.setValue('abc')
    expect(submitBtn.element.disabled).toBe(true)

    await amountInput.setValue('12,50')
    expect(submitBtn.element.disabled).toBe(false)
  })

  it('emittet cancel bei klick auf abbrechen oder close-btn', async () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    await wrapper.find('.btn-cancel').trigger('click')
    expect(wrapper.emitted('cancel')).toBeTruthy()

    await wrapper.find('.close-btn').trigger('click')
    expect(wrapper.emitted('cancel')?.length).toBe(2)
  })

  it('emittet confirm mit geparstem betrag und notiz beim absenden', async () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    const amountInput = wrapper.find<HTMLInputElement>('#expense-amount-input')
    const noteInput = wrapper.find<HTMLInputElement>('#expense-note-input')

    await amountInput.setValue('12,50')
    await noteInput.setValue('Kaffee & Kuchen ')
    await wrapper.find('form').trigger('submit')

    expect(wrapper.emitted('confirm')).toBeTruthy()
    expect(wrapper.emitted('confirm')![0]).toEqual([12.5, 'Kaffee & Kuchen'])
  })

  it('unterstuetzt punkte als dezimaltrennzeichen', async () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    const amountInput = wrapper.find<HTMLInputElement>('#expense-amount-input')

    await amountInput.setValue('8.99')
    await wrapper.find('form').trigger('submit')

    expect(wrapper.emitted('confirm')).toBeTruthy()
    expect(wrapper.emitted('confirm')![0]).toEqual([8.99, ''])
  })

  it('zeigt Lade-Zustand und blockiert Eingaben wenn isSaving=true', async () => {
    const wrapper = mount(Numpad, { props: { visible: true, isSaving: true } })
    expect(wrapper.find('.spinner-inline').exists()).toBe(true)
    expect(wrapper.find('.btn-confirm').text()).toContain('Wird gebucht...')
    expect(wrapper.find<HTMLButtonElement>('.btn-confirm').element.disabled).toBe(true)
    expect(wrapper.find<HTMLInputElement>('#expense-amount-input').element.disabled).toBe(true)
  })

  it('berechnet und zeigt Live Impact Vorschau', async () => {
    const wrapper = mount(Numpad, {
      props: { visible: true, currentBudget: 15.0 },
    })

    // Initialer Stand vor Eingabe
    expect(wrapper.find('.impact-preview').exists()).toBe(true)
    expect(wrapper.find('.impact-preview').classes()).toContain('impact-neutral')
    expect(wrapper.find('.impact-text').text()).toContain('Heute verfügbar: 15,00 €')

    const amountInput = wrapper.find<HTMLInputElement>('#expense-amount-input')
    await amountInput.setValue('8.50')

    expect(wrapper.find('.impact-preview').classes()).toContain('impact-ok')
    expect(wrapper.find('.impact-text').text()).toContain('Verbleibt danach: 6,50 €')

    // Überzogenes Budget
    await amountInput.setValue('20.00')
    expect(wrapper.find('.impact-preview').classes()).toContain('impact-warning')
    expect(wrapper.find('.impact-text').text()).toContain('Überzieht Tagesbudget um 5,00 €')
  })

  it('uebernimmt Notiz bei Klick auf Quick-Chip', async () => {
    const wrapper = mount(Numpad, {
      props: { visible: true },
    })

    const chips = wrapper.findAll('.quick-chip')
    expect(chips.length).toBeGreaterThan(0)

    await chips[0].trigger('click')
    const noteInput = wrapper.find<HTMLInputElement>('#expense-note-input')
    expect(chips[0].text()).toContain(noteInput.element.value)
  })
})

