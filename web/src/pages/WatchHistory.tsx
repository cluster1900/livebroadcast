import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import { message, Card, List, Button, Empty, Spin, Modal, Popconfirm } from 'antd'
import { HistoryOutlined, DeleteOutlined, PlayCircleOutlined, ClockCircleOutlined } from '@ant-design/icons'

interface WatchHistory {
	id: string
	room_id: string
	room_title: string
	streamer_id: string
	streamer_name: string
	cover_url: string
	watch_duration: number
	start_time: string
}

function WatchHistoryPage() {
	const navigate = useNavigate()
	const [loading, setLoading] = useState(true)
	const [history, setHistory] = useState<WatchHistory[]>([])
	const [clearing, setClearing] = useState(false)

	const accessToken = localStorage.getItem('access_token')

	useEffect(() => {
		if (!accessToken) {
			message.warning('请先登录')
			navigate('/login')
			return
		}
		fetchHistory()
	}, [accessToken])

	const fetchHistory = async () => {
		try {
			const response = await axios.get('/api/v1/history/watch', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setHistory(response.data.data)
			}
		} catch (error) {
			message.error('获取观看历史失败')
		} finally {
			setLoading(false)
		}
	}

	const handleDelete = async (id: string) => {
		try {
			const response = await axios.delete(`/api/v1/history/watch/${id}`, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setHistory(history.filter(h => h.id !== id))
				message.success('删除成功')
			}
		} catch (error) {
			message.error('删除失败')
		}
	}

	const handleClearAll = () => {
		Modal.confirm({
			title: '确认清除',
			content: '确定要清除所有观看历史吗？此操作不可恢复。',
			onOk: async () => {
				setClearing(true)
				try {
					const response = await axios.delete('/api/v1/history/watch', {
						headers: { Authorization: `Bearer ${accessToken}` }
					})
					if (response.data.code === 0) {
						setHistory([])
						message.success('已清除所有观看历史')
					}
				} catch (error) {
					message.error('清除失败')
				} finally {
					setClearing(false)
				}
			}
		})
	}

	const formatDuration = (seconds: number) => {
		const hours = Math.floor(seconds / 3600)
		const minutes = Math.floor((seconds % 3600) / 60)
		if (hours > 0) {
			return `${hours}小时${minutes}分钟`
		}
		return `${minutes}分钟`
	}

	const totalWatchTime = history.reduce((sum, h) => sum + h.watch_duration, 0)

	if (loading) {
		return (
			<div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
				<Spin size="large" />
			</div>
		)
	}

	return (
		<div style={{ maxWidth: 900, margin: '0 auto', padding: '20px' }}>
			<Card
				title={
					<span>
						<HistoryOutlined style={{ marginRight: 8 }} />
						观看历史
					</span>
				}
				extra={
					history.length > 0 && (
						<Popconfirm
							title="确定要清除所有观看历史吗？"
							onConfirm={handleClearAll}
							okButtonProps={{ loading: clearing }}
						>
							<Button danger icon={<DeleteOutlined />}>
								清除全部
							</Button>
						</Popconfirm>
					)
				}
			>
				{history.length === 0 ? (
					<Empty
						description="暂无观看记录"
						style={{ padding: 60 }}
					>
						<div style={{ marginBottom: 16 }}>
							去看看精彩的直播吧！
						</div>
						<Button type="primary" onClick={() => navigate('/')}>
							去首页
						</Button>
					</Empty>
				) : (
					<div>
						<div style={{
							background: '#f5f5f5',
							padding: '12px 16px',
							borderRadius: 8,
							marginBottom: 16,
							display: 'flex',
							justifyContent: 'space-between'
						}}>
							<span>共观看 {history.length} 个直播间</span>
							<span>
								<ClockCircleOutlined style={{ marginRight: 4 }} />
								累计观看: {formatDuration(totalWatchTime)}
							</span>
						</div>

						<List
							dataSource={history}
							renderItem={(item) => (
								<List.Item
									style={{
										padding: '12px 0',
										borderBottom: '1px solid #f0f0f0'
									}}
									actions={[
										<Button
											type="primary"
											icon={<PlayCircleOutlined />}
											onClick={() => navigate(`/live/${item.room_id}`)}
										>
											继续观看
										</Button>,
										<Button
											icon={<DeleteOutlined />}
											onClick={() => handleDelete(item.id)}
										/>
									]}
								>
									<List.Item.Meta
										avatar={
											<div
												style={{
													width: 120,
													height: 68,
													borderRadius: 8,
													background: item.cover_url ? `url(${item.cover_url}) center/cover` : '#f0f0f0',
													display: 'flex',
													alignItems: 'center',
													justifyContent: 'center',
													cursor: 'pointer'
												}}
												onClick={() => navigate(`/live/${item.room_id}`)}
											>
												{!item.cover_url && <PlayCircleOutlined style={{ fontSize: 24, color: '#999' }} />}
											</div>
										}
										title={
											<a onClick={() => navigate(`/live/${item.room_id}`)}>
												{item.room_title}
											</a>
										}
										description={
											<div>
												<div>主播: {item.streamer_name}</div>
												<div style={{ color: '#999', fontSize: 12 }}>
													观看时长: {formatDuration(item.watch_duration)} •
													{new Date(item.start_time).toLocaleDateString()}
												</div>
											</div>
										}
									/>
								</List.Item>
							)}
						/>
					</div>
				)}
			</Card>
		</div>
	)
}

export default WatchHistoryPage
