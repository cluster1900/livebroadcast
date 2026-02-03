import { useState, useEffect } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { Avatar, Button, Badge } from 'antd'
import axios from 'axios'
import { BellOutlined, MessageOutlined, GiftOutlined, HistoryOutlined } from '@ant-design/icons'

interface LiveRoom {
  id: string
  title: string
  streamer_id: string
  cover_url: string
  status: string
  peak_online: number
  total_views: number
}

interface UserInfo {
  id: string
  nickname: string
  avatar_url: string
  coin_balance: number
  level: number
}

function Home() {
  const [rooms, setRooms] = useState<LiveRoom[]>([])
  const [userInfo, setUserInfo] = useState<UserInfo | null>(null)
  const [notifUnread, setNotifUnread] = useState(0)
  const [msgUnread, setMsgUnread] = useState(0)
  const navigate = useNavigate()
  const accessToken = localStorage.getItem('access_token')

  useEffect(() => {
    fetchLiveRooms()
    fetchUserInfo()
  }, [])

  const fetchLiveRooms = async () => {
    try {
      const response = await axios.get('/api/v1/live/rooms')
      if (response.data.code === 0) {
        setRooms(response.data.data)
      }
    } catch (error) {
      console.error('Failed to fetch live rooms:', error)
    }
  }

  const fetchUserInfo = async () => {
    const token = localStorage.getItem('access_token')
    if (!token) return

    try {
      const response = await axios.get('/api/v1/user/profile', {
        headers: { Authorization: `Bearer ${token}` }
      })
      if (response.data.code === 0) {
        setUserInfo(response.data.data)
      }
    } catch (error) {
      console.error('Failed to fetch user info:', error)
    }
  }

  const enterRoom = (roomId: string) => {
    navigate(`/live/${roomId}`)
  }

  const handleLogout = () => {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user_id')
    localStorage.removeItem('nickname')
    setUserInfo(null)
  }

  const handleRefreshUnread = async () => {
    if (!accessToken) return
    try {
      const [notifRes, msgRes] = await Promise.all([
        axios.get('/api/v1/notifications/unread-count', { headers: { Authorization: `Bearer ${accessToken}` } }),
        axios.get('/api/v1/messages/unread-count', { headers: { Authorization: `Bearer ${accessToken}` } })
      ])
      setNotifUnread(notifRes.data.data?.count || 0)
      setMsgUnread(msgRes.data.data?.count || 0)
    } catch (e) {
      console.error('获取未读数失败')
    }
  }

  useEffect(() => {
    fetchLiveRooms()
    fetchUserInfo()
    handleRefreshUnread()
    const interval = setInterval(handleRefreshUnread, 30000)
    return () => clearInterval(interval)
  }, [])

  return (
    <div className="home-container">
      <header style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: '20px',
        padding: '16px 24px',
        background: '#fff',
        borderRadius: '8px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
      }}>
        <h2 style={{ margin: 0, color: '#ff6b00' }}>🐯 虎牙直播</h2>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <Button type="text" onClick={() => navigate('/leaderboard')}>排行榜</Button>
          <Button type="text" onClick={() => navigate('/schedules')}>预告</Button>
          {userInfo && <Button type="text" onClick={() => navigate('/streamer')}>主播中心</Button>}
          {userInfo && (
            <Badge count={notifUnread} size="small">
              <Button type="text" icon={<BellOutlined />} onClick={() => navigate('/notifications')}>
                通知
              </Button>
            </Badge>
          )}
          {userInfo && (
            <Badge count={msgUnread} size="small">
              <Button type="text" icon={<MessageOutlined />} onClick={() => navigate('/messages')}>
                私信
              </Button>
            </Badge>
          )}
          {userInfo && (
            <Badge count={0}>
              <Button type="text" icon={<GiftOutlined />} onClick={() => navigate('/inventory')}>
                背包
              </Button>
            </Badge>
          )}
          {userInfo && (
            <Badge count={0}>
              <Button type="text" icon={<HistoryOutlined />} onClick={() => navigate('/history')}>
                历史
              </Button>
            </Badge>
          )}
          {userInfo ? (
            <>
              <Badge count={userInfo.coin_balance} showZero color="#ff6b00" title="虎牙币">
                <Avatar style={{ backgroundColor: '#ff6b00' }} icon="user" src={userInfo.avatar_url} />
              </Badge>
              <span>{userInfo.nickname || '用户'}</span>
              <Link to="/settings">
                <Button type="link">设置</Button>
              </Link>
              <Button onClick={handleLogout} danger>退出</Button>
            </>
          ) : (
            <>
              <Link to="/login">
                <Button type="primary">登录</Button>
              </Link>
              <Link to="/register">
                <Button>注册</Button>
              </Link>
            </>
          )}
        </div>
      </header>

      <h1 style={{ fontSize: '28px', fontWeight: 'bold', marginBottom: '20px' }}>正在直播</h1>
      <div className="live-grid">
        {rooms.map(room => (
          <div key={room.id} className="live-card" onClick={() => enterRoom(room.id)}>
            <div className="live-cover" style={{
              background: room.cover_url ? `url(${room.cover_url}) center/cover` : '#f0f0f0'
            }}>
              {room.status === 'live' && (
                <span style={{
                  position: 'absolute',
                  top: '10px',
                  left: '10px',
                  background: '#ff4d4f',
                  color: 'white',
                  padding: '2px 8px',
                  borderRadius: '4px',
                  fontSize: '12px'
                }}>
                  LIVE
                </span>
              )}
              {room.status === 'running' && (
                <span style={{
                  position: 'absolute',
                  top: '10px',
                  left: '10px',
                  background: '#52c41a',
                  color: 'white',
                  padding: '2px 8px',
                  borderRadius: '4px',
                  fontSize: '12px'
                }}>
                  TV
                </span>
              )}
            </div>
            <div className="live-info">
              <div className="live-title">{room.title}</div>
              <div className="live-meta">
                <span>👁 {room.peak_online}</span>
                <span>🔥 {room.total_views}</span>
              </div>
            </div>
          </div>
        ))}
        {rooms.length === 0 && (
          <div style={{ gridColumn: '1 / -1', textAlign: 'center', padding: '60px', color: '#999' }}>
            <p style={{ fontSize: '18px' }}>暂无直播</p>
            <p style={{ marginTop: '10px' }}>成为主播，开启你的直播之旅！</p>
          </div>
        )}
      </div>
    </div>
  )
}

export default Home
