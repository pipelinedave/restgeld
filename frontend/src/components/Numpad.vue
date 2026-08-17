<template>
  <div v-if="visible" class="numpad-overlay" @click.self="cancel">
    <div class="numpad">
      <div class="numpad-display">
        <span class="numpad-value">{{ display }}</span>
      </div>
      <div class="numpad-grid">
        <button class="numpad-btn" @click="press('7')">7</button>
        <button class="numpad-btn" @click="press('8')">8</button>
        <button class="numpad-btn" @click="press('9')">9</button>
        <button class="numpad-btn" @click="press('4')">4</button>
        <button class="numpad-btn" @click="press('5')">5</button>
        <button class="numpad-btn" @click="press('6')">6</button>
        <button class="numpad-btn" @click="press('1')">1</button>
        <button class="numpad-btn" @click="press('2')">2</button>
        <button class="numpad-btn" @click="press('3')">3</button>
        <button class="numpad-btn" @click="press(',')">,</button>
        <button class="numpad-btn" @click="press('0')">0</button>
        <button class="numpad-btn numpad-btn-del" @click="deleteChar">⌫</button>
        <button class="numpad-btn numpad-btn-cancel" @click="cancel">Abbrechen</button>
        <button class="numpad-btn numpad-btn-confirm" @click="confirm">Bestätigen</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{
  confirm: [value: number]
  cancel: []
}>()

const input = ref('')

const display = computed(() => input.value || '0')

function press(key: string) {
  if (key === ',') {
    if (!input.value.includes(',')) {
      input.value += ','
    }
    return
  }
  input.value += key
}

function deleteChar() {
  input.value = input.value.slice(0, -1)
}

function cancel() {
  input.value = ''
  emit('cancel')
}

function confirm() {
  if (!input.value) return
  const num = parseFloat(input.value.replace(',', '.'))
  if (isNaN(num) || num <= 0) return
  emit('confirm', num)
  input.value = ''
}
</script>

<style scoped>
.numpad-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 100;
}

.numpad {
  width: 100%;
  max-width: 400px;
  background: #112240;
  border-radius: 16px 16px 0 0;
  padding: 16px;
}

.numpad-display {
  text-align: right;
  padding: 8px 4px 16px;
}

.numpad-value {
  font-size: 2rem;
  font-weight: 600;
  color: #ccd6f6;
}

.numpad-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.numpad-btn {
  padding: 16px;
  font-size: 1.3rem;
  border: none;
  border-radius: 8px;
  background: #233554;
  color: #ccd6f6;
  cursor: pointer;
  transition: background 0.15s;
}

.numpad-btn:active {
  background: #1a3a5c;
}

.numpad-btn-del {
  background: #1e3a5f;
  color: #ff6b6b;
}

.numpad-btn-cancel {
  background: #1e3a5f;
  color: #8892b0;
  font-size: 1rem;
}

.numpad-btn-confirm {
  background: #64ffda;
  color: #0a192f;
  font-weight: 700;
  font-size: 1rem;
}

.numpad-grid > :nth-last-child(2) {
  grid-column: 1 / 2;
}

.numpad-grid > :last-child {
  grid-column: 2 / 4;
}
</style>
