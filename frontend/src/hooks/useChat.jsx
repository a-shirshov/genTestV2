import { useState } from 'react'
import { sendPrompt } from '../api/chatApi'

const useChat = () => {
  const [messages, setMessages] = useState([])
  const [isLoading, setIsLoading] = useState(false)
  const [input, setInput] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!input.trim()) return

    try {
      setMessages(prev => [...prev, { text: input, isBot: false }])
      setInput('')
      setIsLoading(true)

      const reader = await sendPrompt(input)
      const decoder = new TextDecoder()
      let botMessage = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) break
        
        const chunk = decoder.decode(value)
        botMessage += chunk
        
        setMessages(prev => {
          const last = prev[prev.length - 1]
          return last?.isBot 
            ? [...prev.slice(0, -1), { text: botMessage, isBot: true }]
            : [...prev, { text: botMessage, isBot: true }]
        })
      }
    } catch (error) {
      console.error('Chat error:', error)
      setMessages(prev => [...prev, { 
        text: `<think>Error: ${error.message}</think>Connection error occurred`,
        isBot: true
      }])
    } finally {
      setIsLoading(false)
    }
  }

  return { messages, isLoading, input, setInput, handleSubmit }
}

export default useChat