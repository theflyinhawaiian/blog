import ReactMarkdown from 'react-markdown'
import rehypeHighlight from 'rehype-highlight'
import 'highlight.js/styles/github.css'

interface Props {
  content: string
  className?: string
}

export function MarkdownContent({ content, className }: Props) {
  return (
    <div className={className}>
      <ReactMarkdown rehypePlugins={[rehypeHighlight]}>{content}</ReactMarkdown>
    </div>
  )
}
