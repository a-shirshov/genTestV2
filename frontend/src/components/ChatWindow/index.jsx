import useChat from '../../hooks/useChat'
import Message from '../Message'
import ChatInput from '../ChatInput'
import styles from './styles.module.css'

const ChatWindow = () => {
  const { messages, isLoading, input, setInput, handleSubmit } = useChat()

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>GenTest Ширшов А.С. ИУ5-44М</h1>
      
      <div className={styles.messagesContainer}>
        {messages.map((msg, i) => (
          <Message key={i} text={msg.text} isBot={msg.isBot} />
        ))}
      </div>

      <ChatInput
        value={input}
        onChange={setInput}
        onSubmit={handleSubmit}
        isLoading={isLoading}
      />
    </div>
  )
}

export default ChatWindow