import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import { message, Card, List, Button, Empty, Spin, Modal, Form, Input, DatePicker, Tag, Tabs, Badge } from 'antd'
import { CalendarOutlined, PlusOutlined, BellOutlined, VideoCameraOutlined } from '@ant-design/icons'
import type { Dayjs } from 'dayjs'

interface Schedule {
	id: string
	streamer_id: string
	streamer_name: string
	title: string
	description: string
	category: string
	cover_url: string
	start_time: string
	status: string
}

function Schedules() {
	const navigate = useNavigate()
	const [loading, setLoading] = useState(true)
	const [upcomingSchedules, setUpcomingSchedules] = useState<Schedule[]>([])
	const [mySchedules, setMySchedules] = useState<Schedule[]>([])
	const [createModalOpen, setCreateModalOpen] = useState(false)
	const [activeTab, setActiveTab] = useState('upcoming')
	const [createForm] = Form.useForm()
	const [creating, setCreating] = useState(false)

	const accessToken = localStorage.getItem('access_token')
	const userId = localStorage.getItem('user_id')

	useEffect(() => {
		if (!accessToken) {
			message.warning('请先登录')
			navigate('/login')
			return
		}
		fetchSchedules()
	}, [accessToken])

	const fetchSchedules = async () => {
		try {
			const [upcomingRes, myRes] = await Promise.all([
				axios.get('/api/v1/extra/schedules/upcoming'),
				axios.get('/api/v1/schedules/my', {
					headers: { Authorization: `Bearer ${accessToken}` }
				})
			])

			if (upcomingRes.data.code === 0) {
				setUpcomingSchedules(upcomingRes.data.data)
			}
			if (myRes.data.code === 0) {
				setMySchedules(myRes.data.data)
			}
		} catch (error) {
			message.error('获取直播预告失败')
		} finally {
			setLoading(false)
		}
	}

	const handleCreateSchedule = async (values: any) => {
		setCreating(true)
		try {
			const response = await axios.post('/api/v1/schedules', {
				...values,
				start_time: values.start_time.toISOString()
			}, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})

			if (response.data.code === 0) {
				message.success('创建成功')
				setCreateModalOpen(false)
				createForm.resetFields()
				fetchSchedules()
			} else {
				message.error(response.data.message || '创建失败')
			}
		} catch (error) {
			message.error('创建失败')
		} finally {
			setCreating(false)
		}
	}

	const handleCancel = async (id: string) => {
		try {
			const response = await axios.post(`/api/v1/schedules/${id}/cancel`, {}, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				message.success('已取消')
				fetchSchedules()
			}
		} catch (error) {
			message.error('取消失败')
		}
	}

	const handleDelete = async (id: string) => {
		try {
			const response = await axios.delete(`/api/v1/schedules/${id}`, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				message.success('删除成功')
				fetchSchedules()
			}
		} catch (error) {
			message.error('删除失败')
		}
	}

	const getCategoryIcon = (category: string) => {
		switch (category) {
			case '游戏': return '🎮'
			case '音乐': return '🎵'
			case '舞蹈': return '💃'
			case '娱乐': return '🎭'
			case '户外': return '🏕️'
			case '体育': return '⚽'
			default: return '📺'
		}
	}

	const getTimeDiff = (time: string) => {
		const diff = new Date(time).getTime() - Date.now()
		const days = Math.floor(diff / (1000 * 60 * 60 * 24))
		const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
		if (days > 0) return `${days}天后`
		if (hours > 0) return `${hours}小时后`
		return '即将开始'
	}

	if (loading) {
		return (
			<div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
				<Spin size="large" />
			</div>
		)
	}

	return (
		<div style={{ maxWidth: 1000, margin: '0 auto', padding: '20px' }}>
			<Card title={<><CalendarOutlined style={{ marginRight: 8 }} />直播预告</>}>
				<Tabs
					activeKey={activeTab}
					onChange={setActiveTab}
					items={[
						{
							key: 'upcoming',
							label: (
								<span>
									即将开播
									{upcomingSchedules.length > 0 && (
										<Badge count={upcomingSchedules.length} style={{ marginLeft: 8 }} />
									)}
								</span>
							),
							children: (
								<ScheduleList
									schedules={upcomingSchedules}
									showActions={false}
									getCategoryIcon={getCategoryIcon}
									getTimeDiff={getTimeDiff}
								/>
							)
						},
						{
							key: 'my',
							label: '我的预告',
							children: (
								<div>
									<Button
										type="primary"
										icon={<PlusOutlined />}
										style={{ marginBottom: 16 }}
										onClick={() => setCreateModalOpen(true)}
									>
										创建预告
									</Button>
									<ScheduleList
										schedules={mySchedules}
										showActions={true}
										isOwner={userId}
										getCategoryIcon={getCategoryIcon}
										getTimeDiff={getTimeDiff}
										onCancel={handleCancel}
										onDelete={handleDelete}
									/>
								</div>
							)
						}
					]}
				/>
			</Card>

			<Modal
				title="创建直播预告"
				open={createModalOpen}
				onCancel={() => setCreateModalOpen(false)}
				footer={null}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreateSchedule}>
					<Form.Item
						name="title"
						label="预告标题"
						rules={[{ required: true, message: '请输入预告标题', min: 2, max: 50 }]}
					>
						<Input placeholder="给你的直播预告起个名字" />
					</Form.Item>
					<Form.Item name="description" label="预告描述">
						<Input.TextArea rows={3} placeholder="描述你的直播内容" />
					</Form.Item>
					<Form.Item name="category" label="直播分类">
						<select style={{ width: '100%', padding: '8px', borderRadius: 4, border: '1px solid #d9d9d9' }}>
							<option value="游戏">游戏</option>
							<option value="音乐">音乐</option>
							<option value="舞蹈">舞蹈</option>
							<option value="娱乐">娱乐</option>
							<option value="户外">户外</option>
							<option value="体育">体育</option>
							<option value="综合">综合</option>
						</select>
					</Form.Item>
					<Form.Item
						name="start_time"
						label="开播时间"
						rules={[{ required: true, message: '请选择开播时间' }]}
					>
						<DatePicker
							showTime
							format="YYYY-MM-DD HH:mm:ss"
							style={{ width: '100%' }}
							disabledDate={(current: Dayjs) => current && current.valueOf() < Date.now()}
						/>
					</Form.Item>
					<Form.Item>
						<Button type="primary" htmlType="submit" loading={creating} block>
							创建预告
						</Button>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	)
}

interface ScheduleListProps {
	schedules: Schedule[]
	showActions?: boolean
	isOwner?: string | null
	getCategoryIcon: (category: string) => string
	getTimeDiff: (time: string) => string
	onCancel?: (id: string) => void
	onDelete?: (id: string) => void
}

function ScheduleList({ schedules, showActions, getCategoryIcon, getTimeDiff, onCancel, onDelete }: ScheduleListProps) {
	if (schedules.length === 0) {
		return <Empty description="暂无直播预告" style={{ padding: 60 }} />
	}

	return (
		<List
			dataSource={schedules}
			renderItem={(item) => (
				<List.Item
					style={{
						padding: '16px',
						borderBottom: '1px solid #f0f0f0',
						background: item.status === 'scheduled' ? 'transparent' : '#fafafa'
					}}
					actions={
						showActions && item.status === 'scheduled' ? [
							<Button size="small" icon={<VideoCameraOutlined />}>
								开播
							</Button>,
							<Button size="small" danger onClick={() => onCancel?.(item.id)}>
								取消
							</Button>
						] : showActions ? [
							<Button size="small" danger onClick={() => onDelete?.(item.id)}>
								删除
							</Button>
						] : [
							<Button type="primary" size="small" icon={<BellOutlined />}>
								提醒我
							</Button>
						]
					}
				>
					<List.Item.Meta
						avatar={
							<div style={{
								width: 120,
								height: 68,
								borderRadius: 8,
								background: item.cover_url ? `url(${item.cover_url}) center/cover` : '#f0f0f0',
								display: 'flex',
								alignItems: 'center',
								justifyContent: 'center',
								fontSize: 24
							}}>
								{!item.cover_url && getCategoryIcon(item.category)}
							</div>
						}
						title={
							<div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
								<span>{item.title}</span>
								<Tag color={item.status === 'scheduled' ? 'blue' : 'default'}>
									{item.status === 'scheduled' ? '预告中' : '已结束'}
								</Tag>
							</div>
						}
						description={
							<div>
								<div>主播: {item.streamer_name}</div>
								<div style={{ color: '#999' }}>
									{getCategoryIcon(item.category)} {item.category} •
									<Tag color="orange" style={{ marginLeft: 4 }}>{getTimeDiff(item.start_time)}</Tag>
								</div>
								<div style={{ color: '#1890ff', fontSize: 12 }}>
									{new Date(item.start_time).toLocaleString()}
								</div>
								{item.description && (
									<div style={{ color: '#666', marginTop: 4, fontSize: 12 }}>
										{item.description}
									</div>
								)}
							</div>
						}
					/>
				</List.Item>
			)}
		/>
	)
}

export default Schedules
