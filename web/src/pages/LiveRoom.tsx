import { useEffect, useState, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import axios from 'axios'
import { message, Spin, Avatar, Input, Button, Drawer, List, Badge, Tooltip, Tag } from 'antd'
import { ShareAltOutlined, UserOutlined, LikeOutlined } from '@ant-design/icons'
import { FLVPlayer } from '../components/VideoPlayer'
import { useCentrifugo } from '../hooks/useCentrifugo'
import type { DanmuMessage, GiftMessage, OnlineCountMessage } from '../components/CentrifugoTypes'
import { motion } from 'framer-motion'

interface RoomInfo {
	id: string
	title: string
	category: string
	cover_url: string
	channel_name: string
	status: string
	streamer_id: string
	streamer_name: string
	start_at: string
	peak_online: number
	total_views: number
	flv_url: string
	hls_url: string
}

interface Danmu {
	id: string
	nickname: string
	level: number
	avatar: string
	content: string
	color: string
}

interface Viewer {
	id: string
	nickname: string
	avatar: string
	level: number
	badge?: string
}

function LiveRoom() {
	const { roomId } = useParams()
	const navigate = useNavigate()
	const [loading, setLoading] = useState(true)
	const [room, setRoom] = useState<RoomInfo | null>(null)
	const [danmuList, setDanmuList] = useState<Danmu[]>([])
	const [onlineCount, setOnlineCount] = useState(0)
	const [giftAnimations, setGiftAnimations] = useState<GiftMessage[]>([])
	const [danmuInput, setDanmuInput] = useState('')
	const [viewersDrawerOpen, setViewersDrawerOpen] = useState(false)
	const [viewers, setViewers] = useState<Viewer[]>([])
	const [playError, setPlayError] = useState(false)
	const danmuRef = useRef<HTMLDivElement>(null)

	const userId = localStorage.getItem('user_id') || ''

	const { connected, connect, disconnect, sendDanmu } = useCentrifugo({
		roomId: roomId || '',
		userId: userId,
		onDanmu: (danmu: DanmuMessage) => {
			setDanmuList(prev => [...prev.slice(-50), {
				id: danmu.data.id,
				nickname: danmu.data.nickname,
				level: danmu.data.level,
				avatar: danmu.data.avatar,
				content: danmu.data.content,
				color: danmu.data.color
			}])
		},
		onGift: (gift: GiftMessage) => {
			setGiftAnimations(prev => [...prev, gift])
			setTimeout(() => {
				setGiftAnimations(prev => prev.slice(1))
			}, 3000)
		},
		onOnlineCount: (data: OnlineCountMessage) => {
			setOnlineCount(data.data.count)
		}
	})

	useEffect(() => {
		fetchRoomInfo()
		return () => {
			disconnect()
		}
	}, [roomId])

	useEffect(() => {
		if (room?.status === 'live' && !connected) {
			connect()
		}
	}, [room, connected])

	useEffect(() => {
		if (danmuRef.current) {
			danmuRef.current.scrollTop = danmuRef.current.scrollHeight
		}
	}, [danmuList])

	const fetchRoomInfo = async () => {
		if (!roomId) return

		try {
			const response = await axios.get(`/api/v1/live/rooms/${roomId}`)
			if (response.data.code === 0) {
				setRoom(response.data.data)
			} else {
				message.error('直播间不存在')
				navigate('/')
			}
		} catch (error) {
			message.error('获取直播间信息失败')
			navigate('/')
		} finally {
			setLoading(false)
		}
	}

	const fetchViewers = async () => {
		if (!roomId || !room) return
		try {
			const response = await axios.get(`/api/v1/social/followers/${room.streamer_id}`)
			if (response.data.code === 0) {
				const followers = response.data.data || []
				const mockViewers: Viewer[] = followers.map((f: any) => ({
					id: f.user_id,
					nickname: f.nickname || '观众',
					avatar: f.avatar_url || '',
					level: Math.floor(Math.random() * 10) + 1,
					badge: f.fan_level >= 5 ? '铁粉' : undefined
				}))
				if (mockViewers.length < 5) {
					mockViewers.push(
						{ id: '1', nickname: '游客001', avatar: '', level: 1 },
						{ id: '2', nickname: '游客002', avatar: '', level: 2 },
						{ id: '3', nickname: '游客003', avatar: '', level: 3 },
					)
				}
				setViewers(mockViewers)
			}
		} catch (error) {
			setViewers([
				{ id: '1', nickname: '游客001', avatar: '', level: 1 },
				{ id: '2', nickname: '游客002', avatar: '', level: 2 },
				{ id: '3', nickname: '游客003', avatar: '', level: 3 },
				{ id: '4', nickname: '游客004', avatar: '', level: 1 },
			])
		}
	}

	const handleShare = () => {
		const shareUrl = window.location.href
		navigator.clipboard.writeText(shareUrl)
		message.success('直播间链接已复制，快分享给好友吧！')
	}

	const handleSendDanmu = async () => {
		if (!danmuInput.trim()) return
		if (!localStorage.getItem('access_token')) {
			message.warning('请先登录')
			navigate('/login')
			return
		}

		const success = await sendDanmu(danmuInput)
		if (success) {
			setDanmuInput('')
		}
	}

	if (loading) {
		return (
			<div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
				<Spin size="large" />
			</div>
		)
	}

	if (!room) {
		return (
			<div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
				<p>直播间不存在</p>
			</div>
		)
	}

	return (
		<div className="room-container">
			<div className="video-section">
				<div style={{ position: 'relative', width: '100%', height: '100%' }}>
					{(room.status === 'live' || room.status === 'running') && !playError ? (
						<FLVPlayer
							src={room.flv_url || room.hls_url}
							poster={room.cover_url}
							onError={() => setPlayError(true)}
						/>
					) : playError ? (
						<div style={{
							display: 'flex',
							flexDirection: 'column',
							alignItems: 'center',
							justifyContent: 'center',
							height: '100%',
							background: '#000'
						}}>
							<div style={{ fontSize: '48px', marginBottom: 16 }}>📺</div>
							<p style={{ fontSize: '24px', color: '#fff' }}>直播已结束</p>
							<p style={{ color: '#999', marginTop: 8 }}>主播: {room.streamer_name}</p>
							<Button type="primary" style={{ marginTop: 16 }} onClick={() => navigate('/')}>
								返回首页
							</Button>
						</div>
					) : (
						<div style={{
							display: 'flex',
							alignItems: 'center',
							justifyContent: 'center',
							height: '100%',
							background: '#000'
						}}>
							<div style={{ textAlign: 'center', color: '#fff' }}>
								<p style={{ fontSize: '24px' }}>直播已结束</p>
								<p style={{ color: '#999' }}>主播: {room.streamer_name}</p>
							</div>
						</div>
					)}

					{giftAnimations.map((gift, index) => (
						<div
							key={`${gift.data.gift.id}-${index}`}
							style={{
								position: 'absolute',
								top: '50%',
								left: '50%',
								transform: 'translate(-50%, -50%)',
								zIndex: 100,
								textAlign: 'center'
							}}
						>
							<motion.div
								initial={{ scale: 0, opacity: 0 }}
								animate={{ scale: 1, opacity: 1 }}
								exit={{ scale: 0, opacity: 0 }}
								transition={{ duration: 0.3 }}
							>
								<div style={{ fontSize: '48px', marginBottom: '10px' }}>
									{gift.data.gift.animation?.includes('lottie') ? '🎁' : '🎁'}
								</div>
								<p style={{ color: '#fff', fontSize: '18px', textShadow: '0 0 10px rgba(0,0,0,0.5)' }}>
									{gift.data.sender.nickname} 送了 {gift.data.gift.name} x{gift.data.count}
								</p>
							</motion.div>
						</div>
					))}
				</div>

				<div className="video-info" style={{
					padding: '15px',
					background: '#1a1a1a',
					color: '#fff',
					display: 'flex',
					justifyContent: 'space-between',
					alignItems: 'center'
				}}>
					<div>
						<h2 style={{ margin: 0, fontSize: '18px' }}>{room.title}</h2>
						<p style={{ margin: '5px 0 0', color: '#999', fontSize: '14px' }}>
							{room.category} • 👁 {room.peak_online}人在看
						</p>
					</div>
					<div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
						<Avatar src={room.cover_url} icon={<UserOutlined />} />
						<span>{room.streamer_name}</span>
						<Badge count={<Tooltip title="当前在线"><LikeOutlined style={{ color: '#52c41a' }} /></Tooltip>} showZero>
							<Button
								type="text"
								icon={<UserOutlined />}
								onClick={() => {
									setViewersDrawerOpen(true)
									fetchViewers()
								}}
							>
								观众列表
							</Button>
						</Badge>
						<Button
							type="text"
							icon={<ShareAltOutlined />}
							onClick={handleShare}
						>
							分享
						</Button>
					</div>
				</div>
			</div>

			<Drawer
				title="👥 观众列表"
				placement="right"
				open={viewersDrawerOpen}
				onClose={() => setViewersDrawerOpen(false)}
				width={300}
			>
				<List
					dataSource={viewers}
					renderItem={(viewer) => (
						<List.Item>
							<List.Item.Meta
								avatar={<Avatar src={viewer.avatar} icon={<UserOutlined />} />}
								title={
									<span>
										{viewer.nickname}
										{viewer.badge && (
											<Tag color="gold" style={{ marginLeft: 8 }}>{viewer.badge}</Tag>
										)}
									</span>
								}
								description={`Lv.${viewer.level}`}
							/>
						</List.Item>
					)}
				/>
			</Drawer>

			<div className="chat-section">
				<div style={{
					padding: '10px 15px',
					borderBottom: '1px solid #e0e0e0',
					fontWeight: 'bold'
				}}>
					弹幕 ({onlineCount}人在线)
				</div>

				<div
					ref={danmuRef}
					className="chat-messages"
					style={{
						height: 'calc(100% - 120px)',
						overflowY: 'auto',
						padding: '10px'
					}}
				>
					{danmuList.map((danmu) => (
						<div key={danmu.id} style={{ marginBottom: '8px' }}>
							<span style={{ color: danmu.color, fontWeight: 'bold' }}>
								{danmu.nickname}
							</span>
							<span style={{ color: '#666', fontSize: '12px', marginLeft: '8px' }}>
								Lv.{danmu.level}
							</span>
							<p style={{ margin: '4px 0 0', color: '#333' }}>
								{danmu.content}
							</p>
						</div>
					))}
					{danmuList.length === 0 && (
						<div style={{ textAlign: 'center', color: '#999', padding: '20px' }}>
							暂无弹幕，快来发第一条吧！
						</div>
					)}
				</div>

				<div className="chat-input" style={{
					padding: '10px',
					borderTop: '1px solid #e0e0e0',
					display: 'flex',
					gap: '8px'
				}}>
					<Input
						placeholder="发送弹幕..."
						value={danmuInput}
						onChange={(e) => setDanmuInput(e.target.value)}
						onPressEnter={handleSendDanmu}
						style={{ flex: 1 }}
					/>
					<button
						onClick={handleSendDanmu}
						style={{
							padding: '8px 16px',
							background: 'linear-gradient(135deg, #ff7e5f 0%, #feb47b 100%)',
							color: 'white',
							border: 'none',
							borderRadius: '4px',
							cursor: 'pointer'
						}}
					>
						发送
					</button>
				</div>
			</div>
		</div>
	)
}

export default LiveRoom
