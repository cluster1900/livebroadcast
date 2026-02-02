import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import Home from './pages/Home'
import LiveRoom from './pages/LiveRoom'
import Login from './pages/Login'
import Register from './pages/Register'
import Settings from './pages/Settings'
import StreamerCenter from './pages/StreamerCenter'
import Leaderboard from './pages/Leaderboard'
import './styles/index.css'

function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/live/:roomId" element={<LiveRoom />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/settings" element={<Settings />} />
          <Route path="/streamer" element={<StreamerCenter />} />
          <Route path="/leaderboard" element={<Leaderboard />} />
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  )
}

export default App
