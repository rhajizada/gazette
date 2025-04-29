import { Github } from 'lucide-react';
import { Link } from "react-router-dom";

export function Footer() {
  return (
    <footer className="w-full py-4 ">
      <div className="mx-auto max-w-5xl px-4 flex items-center justify-between">
        <Link to="/" className="text-lg font-semibold text-gray-800">
          Gazette
        </Link>

        <a
          href="https://github.com/rhajizada/gazette"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center text-gray-600 hover:text-gray-900 transition-colors"
        >
          <Github className="w-6 h-6 mr-2" />
          <span className="text-sm">GitHub</span>
        </a>
      </div>
    </footer>
  )
}
