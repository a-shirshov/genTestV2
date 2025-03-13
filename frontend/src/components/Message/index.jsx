import MarkdownRenderer from '../MarkdownRenderer';
import parseMessage from '../../utils/messageParser';
import styles from './styles.module.css'

const Message = ({ text, isBot }) => {
  const { thinkSection, cleanedText } = parseMessage(text);

  return (
    <div className={`${styles.message} ${isBot ? styles.bot : styles.user}`}>
      {thinkSection && (
        <div className={styles.thinkBubble}>
          <MarkdownRenderer content={thinkSection} />
        </div>
      )}
      <MarkdownRenderer content={cleanedText} />
    </div>
  );
};

export default Message;