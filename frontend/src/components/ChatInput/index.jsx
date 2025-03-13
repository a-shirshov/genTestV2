import styles from './styles.module.css'

const ChatInput = ({ value, onChange, onSubmit, isLoading }) => {
  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      onSubmit(e)
    }
  }

  return (
    <form onSubmit={onSubmit} className={styles.form}>
      <textarea
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="Введите код"
        disabled={isLoading}
        className={styles.input}
      />
      <button 
        type="submit" 
        disabled={isLoading}
        className={styles.button}
      >
        {isLoading ? 'Отправка...' : 'Отправить'}
      </button>
    </form>
  )
}

export default ChatInput