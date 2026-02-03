import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import { message, Card, List, Button, Tabs, Badge, Empty, Spin, Tag } from 'antd'
import { BellOutlined, MailOutlined, GiftOutlined, HeartOutlined, UserAddOutlined, DeleteOutlined } from '@ant-design/icons'

interface Notification {
	id: string
	type: string
	title: string
	content: string
	link: string
	is_read: boolean
	created_at: string
}

function Notifications() {
	const navigate = useNavigate()
	const [loading, setLoading] = useState(true)
	const [notifications, setNotifications] = useState<Notification[]>([])
	const [unreadCount, setUnreadCount] = useState(0)
	const [activeTab, setActiveTab] = useState('all')

	const accessToken = localStorage.getItem('access_token')

	useEffect(() => {
		if (!accessToken) {
			message.warning('请先登录')
			navigate('/login')
			return
		}
		fetchNotifications()
		fetchUnreadCount()
	}, [accessToken])

	const fetchNotifications = async () => {
		try {
			const response = await axios.get('/api/v1/notifications', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setNotifications(response.data.data)
			}
		} catch (error) {
			message.error('获取通知失败')
		} finally {
			setLoading(false)
		}
	}

	const fetchUnreadCount = async () => {
		try {
			const response = await axios.get('/api/v1/notifications/unread-count', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setUnreadCount(response.data.data.count)
			}
		} catch (error) {
			console.error('获取未读数失败')
		}
	}

	const handleMarkAsRead = async (id: string) => {
		try {
			const response = await axios.post(`/api/v1/notifications/${id}/read`, {}, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setNotifications(notifications.map(n =>
					n.id === id ? { ...n, is_read: true } : n
				))
				setUnreadCount(Math.max(0, unreadCount - 1))
			}
		} catch (error) {
			message.error('操作失败')
		}
	}

	const handleMarkAllAsRead = async () => {
		try {
			const response = await axios.post('/api/v1/notifications/read-all', {}, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setNotifications(notifications.map(n => ({ ...n, is_read: true })))
				setUnreadCount(0)
				message.success('已全部标记为已读')
			}
		} catch (error) {
			message.error('操作失败')
		}
	}

	const handleDelete = async (id: string) => {
		try {
			const response = await axios.delete(`/api/v1/notifications/${id}`, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				const deleted = notifications.find(n => n.id === id)
				setNotifications(notifications.filter(n => n.id !== id))
				if (deleted && !deleted.is_read) {
					setUnreadCount(Math.max(0, unreadCount - 1))
				}
				message.success('删除成功')
			}
		} catch (error) {
			message.error('删除失败')
		}
	}

	const getIcon = (type: string) => {
		switch (type) {
			case 'gift': return <GiftOutlined style={{ color: '#ff6b00' }} />
			case 'follow': return <UserAddOutlined style={{ color: '#52c41a' }} />
			case 'like': return <HeartOutlined style={{ color: '#eb2f96' }} />
			case 'system': return <BellOutlined style={{ color: '#1890ff' }} />
			case 'message': return <MailOutlined style={{ color: '#722ed1' }} />
			default: return <BellOutlined />
		}
	}

	const filteredNotifications = activeTab === 'all'
		? notifications
		: activeTab === 'unread'
			? notifications.filter(n => !n.is_read)
			: notifications.filter(n => n.type === activeTab)

	if (loading) {
		return (
			<div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
				<Spin size="large" />
			</div>
		)
	}

	return (
		<div style={{ maxWidth: 800, margin: '0 auto', padding: '20px' }}>
			<Card
				title={
					<span>
						<BellOutlined style={{ marginRight: 8 }} />
						消息通知
						{unreadCount > 0 && (
							<Tag color="red" style={{ marginLeft: 8 }}>{unreadCount}条未读</Tag>
						)}
					</span>
				}
				extra={
					<Button type="link" onClick={handleMarkAllAsRead} disabled={unreadCount === 0}>
						全部已读
					</Button>
				}
			>
				<Tabs
					activeKey={activeTab}
					onChange={setActiveTab}
					items={[
						{
							key: 'all',
							label: '全部',
							children: (
								<NotificationList
									notifications={filteredNotifications}
									onMarkAsRead={handleMarkAsRead}
									onDelete={handleDelete}
									getIcon={getIcon}
								/>
							)
						},
						{
							key: 'unread',
							label: (
								<span>
									未读
									{unreadCount > 0 && <Badge count={unreadCount} style={{ marginLeft: 8 }} />}
								</span>
							),
							children: (
								<NotificationList
									notifications={filteredNotifications}
									onMarkAsRead={handleMarkAsRead}
									onDelete={handleDelete}
									getIcon={getIcon}
								/>
							)
						},
						{
							key: 'system',
							label: '系统',
							children: (
								<NotificationList
									notifications={filteredNotifications}
									onMarkAsRead={handleMarkAsRead}
									onDelete={handleDelete}
									getIcon={getIcon}
								/>
							)
						},
						{
							key: 'gift',
							label: '礼物',
							children: (
								<NotificationList
									notifications={filteredNotifications}
									onMarkAsRead={handleMarkAsRead}
									onDelete={handleDelete}
									getIcon={getIcon}
								/>
							)
						},
						{
							key: 'social',
							label: '互动',
							children: (
								<NotificationList
									notifications={filteredNotifications}
									onMarkAsRead={handleMarkAsRead}
									onDelete={handleDelete}
									getIcon={getIcon}
								/>
							)
						}
					]}
				/>
			</Card>
		</div>
	)
}

function NotificationList({ notifications, onMarkAsRead, onDelete, getIcon }: {
	notifications: Notification[]
	onMarkAsRead: (id: string) => void
	onDelete: (id: string) => void
	getIcon: (type: string) => React.ReactNode
}) {
	if (notifications.length === 0) {
		return <Empty description="暂无通知" style={{ padding: 40 }} />
	}

	return (
		<List
			dataSource={notifications}
			renderItem={(item) => (
				<List.Item
					style={{
						background: item.is_read ? 'transparent' : '#f6ffed',
						padding: '12px 16px',
						marginBottom: 8,
						borderRadius: 8
					}}
					actions={[
						!item.is_read && (
							<Button type="link" size="small" onClick={() => onMarkAsRead(item.id)}>
								标记已读
							</Button>
						),
						<Button
							type="text"
							danger
							size="small"
							icon={<DeleteOutlined />}
							onClick={() => onDelete(item.id)}
						/>
					].filter(Boolean)}
				>
					<List.Item.Meta
						avatar={
							<div style={{
								fontSize: 24,
								width: 40,
								height: 40,
								display: 'flex',
								alignItems: 'center',
								justifyContent: 'center',
								background: '#f5f5f5',
								borderRadius: '50%'
							}}>
								{getIcon(item.type)}
							</div>
						}
						title={
							<span>
								{item.title}
								{item.is_read && (
									<Tag color="default" style={{ marginLeft: 8 }}>已读</Tag>
								)}
							</span>
						}
						description={
							<div>
								<div style={{ color: '#666', marginBottom: 4 }}>{item.content}</div>
								<div style={{ color: '#999', fontSize: 12 }}>
									{new Date(item.created_at).toLocaleString()}
								</div>
							</div>
						}
					/>
				</List.Item>
			)}
		/>
	)
}

export default Notifications
