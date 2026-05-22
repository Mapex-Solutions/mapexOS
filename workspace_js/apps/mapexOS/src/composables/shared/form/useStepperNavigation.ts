import { onMounted, onUnmounted } from 'vue'
import type { Ref } from 'vue'

interface UseStepperNavigationOptions {
  currentStep: Ref<number>
  totalSteps: number
  changeStep: (step: number) => void
}

export function useStepperNavigation({
  currentStep,
  totalSteps,
  changeStep
}: UseStepperNavigationOptions) {
  function handleKey(event: KeyboardEvent) {
    // ignore if focused in form inputs
    if (
      event.target instanceof HTMLInputElement ||
      event.target instanceof HTMLTextAreaElement ||
      event.target instanceof HTMLSelectElement
    ) {
      return
    }

    switch (event.key) {
      case 'ArrowLeft':
        if (currentStep.value > 1) {
          changeStep(currentStep.value - 1)
        }
        event.preventDefault()
        break

      case 'ArrowRight':
        if (currentStep.value < totalSteps) {
          changeStep(currentStep.value + 1)
        }
        event.preventDefault()
        break
    }
  }

  onMounted(() => {
    document.addEventListener('keydown', handleKey)
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKey)
  })
}
