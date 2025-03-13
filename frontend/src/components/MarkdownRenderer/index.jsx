import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { atomDark } from 'react-syntax-highlighter/dist/esm/styles/prism'
import { useState } from 'react'
import ReactMarkdown from 'react-markdown';
import styles from './styles.module.css'

const MarkdownRenderer = ({ content }) => {
  const [copied, setCopied] = useState(false)

  const handleCopy = (code) => {
    navigator.clipboard.writeText(code)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className={styles.container}>
      <ReactMarkdown
        children={content}
        components={{
          code({ node, inline, className, children, ...props }) {
            const match = /language-(\w+)/.exec(className || '')
            const code = String(children).replace(/\n$/, '')
            
            return !inline ? (
              <div className={styles.codeBlock}>
                <div className={styles.codeHeader}>
                  <span>{match ? match[1] : 'Code'}</span>
                  <button 
                    onClick={() => handleCopy(code)}
                    className={styles.copyButton}
                  >
                    {copied ? '✓ Copied' : 'Copy'}
                  </button>
                </div>
                <SyntaxHighlighter
                  children={code}
                  style={atomDark}
                  language={match?.[1] || 'text'}
                  PreTag="div"
                  showLineNumbers
                  {...props}
                />
              </div>
            ) : (
              <code className={styles.inlineCode} {...props}>
                {children}
              </code>
            )
          },
          a({ href, children }) {
            return <a href={href} target="_blank" rel="noopener noreferrer" className={styles.link}>{children}</a>
          }
        }}
      />
    </div>
  )
}

export default MarkdownRenderer