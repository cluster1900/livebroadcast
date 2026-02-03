import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import { message, Card, Layout, Menu, Table, Button, Tag, Modal, Form, Input, Statistic, Row, Col, Tabs, Badge, Popconfirm, Select } from 'antd'
import { DashboardOutlined, UserOutlined, VideoCameraOutlined, GiftOutlined, WarningOutlined, SettingOutlined, TeamOutlined, BarChartOutlined } from '@ant-design/icons'

const { Header, Sider, Content } = Layout

interface DashboardStats {
	total_users: number
	total_streamers: number
	total_rooms: number
	live_rooms: number
	total_revenue: number
	pending_reports: number
	new_users_today: number
}

interface User {
	id: string
	username: string
	nickname: string
	level: number
	coin_balance: number
	status: string
	created_at: string
}

interface Room {
	id: string
	title: string
	streamer_id: string
	streamer: string
	status: string
	peak_online: number
	total_views: number
	created_at: string
}

interface Gift {
	id: number
	name: string
	coin_price: number
	icon_url: string
	category: string
	sort_order: number
	is_active: boolean
}

interface SensitiveWord {
	id: number
	word: string
	type: string
	severity: string
	is_active: boolean
}

interface Report {
	id: string
	reporter: string
	reported: string
	type: string
	reason: string
	status: string
	created_at: string
}

function Admin() {
	const navigate = useNavigate()
	const [collapsed, setCollapsed] = useState(false)
	const [loading, setLoading] = useState(true)
	const [stats, setStats] = useState<DashboardStats | null>(null)
	const [users, setUsers] = useState<User[]>([])
	const [rooms, setRooms] = useState<Room[]>([])
	const [gifts, setGifts] = useState<Gift[]>([])
	const [words, setWords] = useState<SensitiveWord[]>([])
	const [reports, setReports] = useState<Report[]>([])
	const [activeMenu, setActiveMenu] = useState('dashboard')
	const [giftModalOpen, setGiftModalOpen] = useState(false)
	const [wordModalOpen, setWordModalOpen] = useState(false)
	const [editingGift, setEditingGift] = useState<Gift | null>(null)
	const [handleModalOpen, setHandleModalOpen] = useState(false)
	const [selectedReport, setSelectedReport] = useState<Report | null>(null)
	const [handleNote, setHandleNote] = useState('')
	const [giftForm] = Form.useForm()
	const [wordForm] = Form.useForm()

	const accessToken = localStorage.getItem('access_token')
	const userRole = localStorage.getItem('user_role')

	useEffect(() => {
		if (!accessToken || userRole !== 'admin') {
			message.warning('请使用管理员账号登录')
			navigate('/login')
			return
		}
		fetchDashboard()
		fetchUsers()
		fetchRooms()
		fetchGifts()
		fetchWords()
		fetchReports()
	}, [accessToken, userRole])

	const fetchDashboard = async () => {
		try {
			const response = await axios.get('/api/v1/admin/dashboard', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setStats(response.data.data)
			}
		} catch (error) {
			console.error('获取统计失败')
		} finally {
			setLoading(false)
		}
	}

	const fetchUsers = async () => {
		try {
			const response = await axios.get('/api/v1/admin/users', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setUsers(response.data.data)
			}
		} catch (error) {
			console.error('获取用户失败')
		}
	}

	const fetchRooms = async () => {
		try {
			const response = await axios.get('/api/v1/admin/rooms', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setRooms(response.data.data)
			}
		} catch (error) {
			console.error('获取直播间失败')
		}
	}

	const fetchGifts = async () => {
		try {
			const response = await axios.get('/api/v1/admin/gifts', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setGifts(response.data.data)
			}
		} catch (error) {
			console.error('获取礼物失败')
		}
	}

	const fetchWords = async () => {
		try {
			const response = await axios.get('/api/v1/admin/sensitive-words', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setWords(response.data.data)
			}
		} catch (error) {
			console.error('获取敏感词失败')
		}
	}

	const fetchReports = async () => {
		try {
			const response = await axios.get('/api/v1/admin/reports/pending', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setReports(response.data.data)
			}
		} catch (error) {
			console.error('获取举报失败')
		}
	}

	const handleBanUser = async (userId: string) => {
		try {
			const response = await axios.post(`/api/v1/admin/users/${userId}/ban`, {}, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				message.success('已封禁用户')
				fetchUsers()
			}
		} catch (error) {
			message.error('操作失败')
		}
	}

	const handleUnbanUser = async (userId: string) => {
		try {
			const response = await axios.post(`/api/v1/admin/users/${userId}/unban`, {}, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				message.success('已解封用户')
				fetchUsers()
			}
		} catch (error) {
			message.error('操作失败')
		}
	}

	const handleBanRoom = async (roomId: string, reason: string) => {
		try {
			const response = await axios.post(`/api/v1/admin/rooms/${roomId}/ban?reason=${reason}`, {}, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				message.success('已封禁直播间')
				fetchRooms()
			}
		} catch (error) {
			message.error('操作失败')
		}
	}

	const handleReport = async (status: string) => {
		if (!selectedReport) return
		try {
			const response = await axios.post(`/api/v1/admin/reports/${selectedReport.id}/handle`, {
				status,
				handle_note: handleNote
			}, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				message.success('处理完成')
				setHandleModalOpen(false)
				setSelectedReport(null)
				setHandleNote('')
				fetchReports()
			}
		} catch (error) {
			message.error('处理失败')
		}
	}

	const handleSaveGift = async (values: any) => {
		try {
			if (editingGift) {
				await axios.put(`/api/v1/admin/gifts/${editingGift.id}`, values, {
					headers: { Authorization: `Bearer ${accessToken}` }
				})
				message.success('更新成功')
			} else {
				await axios.post('/api/v1/admin/gifts', values, {
					headers: { Authorization: `Bearer ${accessToken}` }
				})
				message.success('创建成功')
			}
			setGiftModalOpen(false)
			setEditingGift(null)
			giftForm.resetFields()
			fetchGifts()
		} catch (error) {
			message.error('操作失败')
		}
	}

	const handleDeleteGift = async (id: number) => {
		try {
			await axios.delete(`/api/v1/admin/gifts/${id}`, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			message.success('删除成功')
			fetchGifts()
		} catch (error) {
			message.error('删除失败')
		}
	}

	const handleSaveWord = async (values: any) => {
		try {
			await axios.post('/api/v1/admin/sensitive-words', values, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			message.success('添加成功')
			setWordModalOpen(false)
			wordForm.resetFields()
			fetchWords()
		} catch (error) {
			message.error('添加失败')
		}
	}

	const handleDeleteWord = async (id: number) => {
		try {
			await axios.delete(`/api/v1/admin/sensitive-words/${id}`, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			message.success('删除成功')
			fetchWords()
		} catch (error) {
			message.error('删除失败')
		}
	}

	const menuItems = [
		{ key: 'dashboard', icon: <DashboardOutlined />, label: '数据概览' },
		{ key: 'users', icon: <TeamOutlined />, label: '用户管理' },
		{ key: 'rooms', icon: <VideoCameraOutlined />, label: '直播间管理' },
		{ key: 'gifts', icon: <GiftOutlined />, label: '礼物管理' },
		{ key: 'reports', icon: <WarningOutlined />, label: '举报处理' },
		{ key: 'words', icon: <WarningOutlined />, label: '敏感词管理' },
		{ key: 'settings', icon: <SettingOutlined />, label: '系统配置' },
	]

	const userColumns = [
		{ title: '用户名', dataIndex: 'username', key: 'username' },
		{ title: '昵称', dataIndex: 'nickname', key: 'nickname' },
		{ title: '等级', dataIndex: 'level', key: 'level', render: (l: number) => `Lv.${l}` },
		{ title: '余额', dataIndex: 'coin_balance', key: 'balance', render: (b: number) => `${b} 币` },
		{
			title: '状态',
			dataIndex: 'status',
			key: 'status',
			render: (s: string) => <Tag color={s === 'active' ? 'green' : 'red'}>{s === 'active' ? '正常' : '封禁'}</Tag>
		},
		{
			title: '操作',
			key: 'action',
			render: (_: any, record: User) => (
				<div>
					{record.status === 'active' ? (
						<Popconfirm title="确定封禁此用户？" onConfirm={() => handleBanUser(record.id)}>
							<Button type="link" danger size="small">封禁</Button>
						</Popconfirm>
					) : (
						<Button type="link" size="small" onClick={() => handleUnbanUser(record.id)}>解封</Button>
					)}
				</div>
			)
		}
	]

	const roomColumns = [
		{ title: '标题', dataIndex: 'title', key: 'title', ellipsis: true },
		{ title: '主播', dataIndex: 'streamer', key: 'streamer' },
		{
			title: '状态',
			dataIndex: 'status',
			key: 'status',
			render: (s: string) => (
				<Tag color={s === 'live' ? 'green' : s === 'banned' ? 'red' : 'default'}>
					{s === 'live' ? '直播中' : s === 'banned' ? '已封禁' : '已结束'}
				</Tag>
			)
		},
		{ title: '观看峰值', dataIndex: 'peak_online', key: 'peak' },
		{ title: '总观看', dataIndex: 'total_views', key: 'views' },
		{
			title: '操作',
			key: 'action',
			render: (_: any, record: Room) => (
				record.status !== 'banned' && (
					<Popconfirm title="确定封禁此直播间？" onConfirm={() => handleBanRoom(record.id, '违规内容')}>
						<Button type="link" danger size="small">封禁</Button>
					</Popconfirm>
				)
			)
		}
	]

	const giftColumns = [
		{ title: '名称', dataIndex: 'name', key: 'name' },
		{ title: '价格', dataIndex: 'coin_price', key: 'price', render: (p: number) => `${p} 币` },
		{ title: '分类', dataIndex: 'category', key: 'category' },
		{ title: '排序', dataIndex: 'sort_order', key: 'order' },
		{
			title: '状态',
			dataIndex: 'is_active',
			key: 'active',
			render: (a: boolean) => <Tag color={a ? 'green' : 'red'}>{a ? '启用' : '禁用'}</Tag>
		},
		{
			title: '操作',
			key: 'action',
			render: (_: any, record: Gift) => (
				<div>
					<Button type="link" size="small" onClick={() => {
						setEditingGift(record)
						giftForm.setFieldsValue(record)
						setGiftModalOpen(true)
					}}>编辑</Button>
					<Popconfirm title="确定删除？" onConfirm={() => handleDeleteGift(record.id)}>
						<Button type="link" danger size="small">删除</Button>
					</Popconfirm>
				</div>
			)
		}
	]

	const wordColumns = [
		{ title: '敏感词', dataIndex: 'word', key: 'word' },
		{
			title: '类型',
			dataIndex: 'type',
			key: 'type',
			render: (t: string) => <Tag>{t === 'blacklist' ? '黑名单' : '白名单'}</Tag>
		},
		{
			title: '等级',
			dataIndex: 'severity',
			key: 'severity',
			render: (s: string) => <Tag color={s === 'high' ? 'red' : s === 'medium' ? 'orange' : 'blue'}>{s}</Tag>
		},
		{
			title: '操作',
			key: 'action',
			render: (_: any, record: SensitiveWord) => (
				<Popconfirm title="确定删除？" onConfirm={() => handleDeleteWord(record.id)}>
					<Button type="link" danger size="small">删除</Button>
				</Popconfirm>
			)
		}
	]

	const reportColumns = [
		{ title: '举报人', dataIndex: 'reporter', key: 'reporter' },
		{ title: '被举报人', dataIndex: 'reported', key: 'reported' },
		{ title: '类型', dataIndex: 'type', key: 'type' },
		{ title: '原因', dataIndex: 'reason', key: 'reason', ellipsis: true },
		{ title: '时间', dataIndex: 'created_at', key: 'time', render: (t: string) => new Date(t).toLocaleString() },
		{
			title: '操作',
			key: 'action',
			render: (_: any, record: Report) => (
				<Button type="primary" size="small" onClick={() => {
					setSelectedReport(record)
					setHandleModalOpen(true)
				}}>处理</Button>
			)
		}
	]

	const renderContent = () => {
		switch (activeMenu) {
			case 'dashboard':
				return (
					<Row gutter={16}>
						<Col span={6}><Card><Statistic title="用户总数" value={stats?.total_users || 0} /></Card></Col>
						<Col span={6}><Card><Statistic title="主播数" value={stats?.total_streamers || 0} /></Card></Col>
						<Col span={6}><Card><Statistic title="直播间数" value={stats?.total_rooms || 0} /></Card></Col>
						<Col span={6}><Card><Statistic title="直播中" value={stats?.live_rooms || 0} valueStyle={{ color: '#52c41a' }} /></Card></Col>
						<Col span={6}><Card><Statistic title="总收入" value={stats?.total_revenue || 0} prefix="¥" /></Card></Col>
						<Col span={6}><Card><Statistic title="今日新增" value={stats?.new_users_today || 0} /></Card></Col>
						<Col span={6}>
							<Card>
								<Statistic
									title="待处理举报"
									value={stats?.pending_reports || 0}
									valueStyle={{ color: stats?.pending_reports ? '#ff4d4f' : '#52c41a' }}
									prefix={<Badge count={stats?.pending_reports} />}
								/>
							</Card>
						</Col>
					</Row>
				)
			case 'users':
				return <Table dataSource={users} columns={userColumns} rowKey="id" />
			case 'rooms':
				return <Table dataSource={rooms} columns={roomColumns} rowKey="id" />
			case 'gifts':
				return (
					<div>
						<Button type="primary" style={{ marginBottom: 16 }} onClick={() => {
							setEditingGift(null)
							giftForm.resetFields()
							setGiftModalOpen(true)
						}}>添加礼物</Button>
						<Table dataSource={gifts} columns={giftColumns} rowKey="id" />
					</div>
				)
			case 'reports':
				return <Table dataSource={reports} columns={reportColumns} rowKey="id" />
			case 'words':
				return (
					<div>
						<Button type="primary" style={{ marginBottom: 16 }} onClick={() => setWordModalOpen(true)}>添加敏感词</Button>
						<Table dataSource={words} columns={wordColumns} rowKey="id" />
					</div>
				)
			case 'settings':
				return <Card title="系统配置"><p>系统配置功能开发中...</p></Card>
			default:
				return null
		}
	}

	return (
		<Layout style={{ minHeight: '100vh' }}>
			<Sider collapsible collapsed={collapsed} onCollapse={setCollapsed}>
				<div style={{ height: 64, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#fff', fontSize: 18 }}>
					{collapsed ? '🐯' : '🐯 虎牙管理后台'}
				</div>
				<Menu theme="dark" mode="inline" selectedKeys={[activeMenu]} items={menuItems} onClick={(e) => setActiveMenu(e.key)} />
			</Sider>
			<Layout>
				<Header style={{ padding: '0 24px', background: '#fff', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
					<h3>{menuItems.find(m => m.key === activeMenu)?.label}</h3>
					<Button onClick={() => navigate('/')}>返回前台</Button>
				</Header>
				<Content style={{ margin: 16, padding: 24, background: '#fff' }}>
					{renderContent()}
				</Content>
			</Layout>

			<Modal title={editingGift ? '编辑礼物' : '添加礼物'} open={giftModalOpen} onCancel={() => setGiftModalOpen(false)} footer={null}>
				<Form form={giftForm} layout="vertical" onFinish={handleSaveGift}>
					<Form.Item name="name" label="名称" rules={[{ required: true }]}><Input /></Form.Item>
					<Form.Item name="coin_price" label="价格" rules={[{ required: true }]}><Input type="number" /></Form.Item>
					<Form.Item name="icon_url" label="图标URL"><Input /></Form.Item>
					<Form.Item name="category" label="分类">
						<Select options={[{ value: 'normal', label: '普通' }, { value: 'vip', label: 'VIP' }, { value: 'special', label: 'Special' }]} />
					</Form.Item>
					<Form.Item name="sort_order" label="排序"><Input type="number" /></Form.Item>
					<Form.Item>
						<Button type="primary" htmlType="submit" block>保存</Button>
					</Form.Item>
				</Form>
			</Modal>

			<Modal title="添加敏感词" open={wordModalOpen} onCancel={() => setWordModalOpen(false)} footer={null}>
				<Form form={wordForm} layout="vertical" onFinish={handleSaveWord}>
					<Form.Item name="word" label="敏感词" rules={[{ required: true }]}><Input /></Form.Item>
					<Form.Item name="type" label="类型" initialValue="blacklist">
						<Select options={[{ value: 'blacklist', label: '黑名单' }, { value: 'whitelist', label: '白名单' }]} />
					</Form.Item>
					<Form.Item name="severity" label="等级" initialValue="medium">
						<Select options={[{ value: 'low', label: 'Low' }, { value: 'medium', label: 'Medium' }, { value: 'high', label: 'High' }]} />
					</Form.Item>
					<Form.Item>
						<Button type="primary" htmlType="submit" block>添加</Button>
					</Form.Item>
				</Form>
			</Modal>

			<Modal title="处理举报" open={handleModalOpen} onCancel={() => setHandleModalOpen(false)} footer={null}>
				<p><strong>举报人:</strong> {selectedReport?.reporter}</p>
				<p><strong>被举报人:</strong> {selectedReport?.reported}</p>
				<p><strong>类型:</strong> {selectedReport?.type}</p>
				<p><strong>原因:</strong> {selectedReport?.reason}</p>
				<Input.TextArea rows={3} placeholder="处理备注" value={handleNote} onChange={(e) => setHandleNote(e.target.value)} style={{ marginBottom: 16 }} />
				<div style={{ display: 'flex', gap: 8 }}>
					<Button type="primary" onClick={() => handleReport('resolved')}>通过</Button>
					<Button danger onClick={() => handleReport('dismissed')}>驳回</Button>
				</div>
			</Modal>
		</Layout>
	)
}

export default Admin
