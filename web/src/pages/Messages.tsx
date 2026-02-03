import { useState, useEffect, useRef } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import axios from 'axios'
import { message, Card, List, Avatar, Input, Button, Badge, Spin, Empty } from 'antd'
import { SendOutlined, UserOutlined } from '@ant-design/icons'
import { motion, AnimatePresence } from 'framer-motion'

interface Conversation {
	user_id: string
	nickname: string
	avatar_url: string
	last_message: string
	last_time: string
	unread_count: number
}

interface Message {
	id: string
	sender_id: string
	receiver_id: string
	content: string
	is_read: boolean
	created_at: string
}

function Messages() {
	const { userId: targetUserId } = useParams()
	const navigate = useNavigate()
	const [loading, setLoading] = useState(true)
	const [conversations, setConversations] = useState<Conversation[]>([])
	const [selectedUser, setSelectedUser] = useState<string | null>(targetUserId || null)
	const [messages, setMessages] = useState<Message[]>([])
	const [newMessage, setNewMessage] = useState('')
	const [sending, setSending] = useState(false)
	const [unreadCount, setUnreadCount] = useState(0)
	const messagesEndRef = useRef<HTMLDivElement>(null)

	const accessToken = localStorage.getItem('access_token')
	const currentUserId = localStorage.getItem('user_id')

	useEffect(() => {
		if (!accessToken) {
			message.warning('请先登录')
			navigate('/login')
			return
		}
		fetchConversations()
		fetchUnreadCount()
	}, [accessToken])

	useEffect(() => {
		if (selectedUser) {
			fetchMessages(selectedUser)
		}
	}, [selectedUser])

	useEffect(() => {
		scrollToBottom()
	}, [messages])

	const fetchConversations = async () => {
		try {
			const response = await axios.get('/api/v1/messages/conversations', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setConversations(response.data.data)
				if (!selectedUser && response.data.data.length > 0) {
					setSelectedUser(response.data.data[0].user_id)
				}
			}
		} catch (error) {
			message.error('获取会话列表失败')
		} finally {
			setLoading(false)
		}
	}

	const fetchMessages = async (userId: string) => {
		try {
			const response = await axios.get(`/api/v1/messages/with/${userId}`, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setMessages(response.data.data.reverse())
			}
		} catch (error) {
			message.error('获取消息失败')
		}
	}

	const fetchUnreadCount = async () => {
		try {
			const response = await axios.get('/api/v1/messages/unread-count', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setUnreadCount(response.data.data.count)
			}
		} catch (error) {
			console.error('获取未读数失败')
		}
	}

	const handleSend = async () => {
		if (!newMessage.trim() || !selectedUser) return

		setSending(true)
		try {
			const response = await axios.post('/api/v1/messages/send', {
				receiver_id: selectedUser,
				content: newMessage.trim()
			}, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})

			if (response.data.code === 0) {
				const newMsg: Message = {
					id: response.data.data.id,
					sender_id: currentUserId!,
					receiver_id: selectedUser,
					content: newMessage.trim(),
					is_read: false,
					created_at: new Date().toISOString()
				}
				setMessages([...messages, newMsg])
				setNewMessage('')
			} else {
				message.error(response.data.message || '发送失败')
			}
		} catch (error) {
			message.error('发送失败')
		} finally {
			setSending(false)
		}
	}

	const handleSelectUser = (userId: string) => {
		setSelectedUser(userId)
		navigate(`/messages/${userId}`)
	}

	const scrollToBottom = () => {
		messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
	}

	if (loading) {
		return (
			<div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
				<Spin size="large" />
			</div>
		)
	}

	return (
		<div style={{ maxWidth: 1000, margin: '0 auto', padding: '20px', height: 'calc(100vh - 40px)' }}>
			<Card style={{ height: '100%' }} bodyStyle={{ display: 'flex', height: '100%', padding: 0 }}>
				<div style={{ width: 300, borderRight: '1px solid #f0f0f0', display: 'flex', flexDirection: 'column' }}>
					<div style={{ padding: '16px', borderBottom: '1px solid #f0f0f0' }}>
						<h3>
							私信
							{unreadCount > 0 && (
								<Badge count={unreadCount} style={{ marginLeft: 8 }} />
							)}
						</h3>
					</div>
					<div style={{ flex: 1, overflow: 'auto' }}>
						{conversations.length === 0 ? (
							<Empty description="暂无私信" style={{ padding: 40 }} />
						) : (
							<List
								dataSource={conversations}
								renderItem={(conv) => (
									<List.Item
										style={{
											padding: '12px 16px',
											cursor: 'pointer',
											background: selectedUser === conv.user_id ? '#e6f7ff' : 'transparent'
										}}
										onClick={() => handleSelectUser(conv.user_id)}
									>
										<List.Item.Meta
											avatar={
												<Badge count={conv.unread_count} size="small">
													<Avatar src={conv.avatar_url} icon={<UserOutlined />} />
												</Badge>
											}
											title={
												<span style={{ fontWeight: conv.unread_count > 0 ? 'bold' : 'normal' }}>
													{conv.nickname}
												</span>
											}
											description={
												<div>
													<div style={{
														color: '#999',
														overflow: 'hidden',
														textOverflow: 'ellipsis',
														whiteSpace: 'nowrap',
														maxWidth: 180
													}}>
														{conv.last_message}
													</div>
													<div style={{ color: '#ccc', fontSize: 12 }}>
														{new Date(conv.last_time).toLocaleDateString()}
													</div>
												</div>
											}
										/>
									</List.Item>
								)}
							/>
						)}
					</div>
				</div>

				<div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
					{selectedUser ? (
						<>
							<div style={{
								padding: '12px 16px',
								borderBottom: '1px solid #f0f0f0',
								display: 'flex',
								alignItems: 'center'
							}}>
								<Avatar src={conversations.find(c => c.user_id === selectedUser)?.avatar_url} icon={<UserOutlined />} />
								<span style={{ marginLeft: 12, fontWeight: 'bold' }}>
									{conversations.find(c => c.user_id === selectedUser)?.nickname}
								</span>
							</div>

							<div style={{ flex: 1, overflow: 'auto', padding: '16px' }}>
								{messages.length === 0 ? (
									<Empty description="暂无消息，开始聊天吧" />
								) : (
									<AnimatePresence>
										{messages.map((msg) => {
											const isOwn = msg.sender_id === currentUserId
											return (
												<motion.div
													key={msg.id}
													initial={{ opacity: 0, y: 20 }}
													animate={{ opacity: 1, y: 0 }}
													exit={{ opacity: 0 }}
													style={{
														display: 'flex',
														justifyContent: isOwn ? 'flex-end' : 'flex-start',
														marginBottom: 16
													}}
												>
													<div style={{
														maxWidth: '70%',
														padding: '10px 14px',
														borderRadius: 12,
														background: isOwn ? '#ff6b00' : '#f0f0f0',
														color: isOwn ? '#fff' : '#333'
													}}>
														{msg.content}
													</div>
												</motion.div>
											)
										})}
									</AnimatePresence>
								)}
								<div ref={messagesEndRef} />
							</div>

							<div style={{
								padding: '12px 16px',
								borderTop: '1px solid #f0f0f0',
								display: 'flex',
								gap: 8
							}}>
								<Input
									placeholder="输入消息..."
									value={newMessage}
									onChange={(e) => setNewMessage(e.target.value)}
									onPressEnter={handleSend}
									style={{ flex: 1 }}
								/>
								<Button
									type="primary"
									icon={<SendOutlined />}
									onClick={handleSend}
									loading={sending}
									disabled={!newMessage.trim()}
								>
									发送
								</Button>
							</div>
						</>
					) : (
						<div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
							<Empty description="选择一个对话开始聊天" />
						</div>
					)}
				</div>
			</Card>
		</div>
	)
}

export default Messages
