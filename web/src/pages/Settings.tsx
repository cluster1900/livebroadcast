import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import { message, Card, Form, Input, Button, Upload, Avatar, Tabs, List, Spin } from 'antd'

interface UserProfile {
	id: string
	username: string
	nickname: string
	avatar_url: string
	email: string
	phone: string
	level: number
	exp: number
	coin_balance: number
}

interface Transaction {
	id: number
	amount: number
	balance_after: number
	type: string
	description: string
	created_at: string
}

function Settings() {
	const navigate = useNavigate()
	const [loading, setLoading] = useState(true)
	const [profile, setProfile] = useState<UserProfile | null>(null)
	const [transactions, setTransactions] = useState<Transaction[]>([])
	const [editForm] = Form.useForm()
	const [rechargeForm] = Form.useForm()
	const [activeTab, setActiveTab] = useState('profile')

	const userId = localStorage.getItem('user_id')
	const accessToken = localStorage.getItem('access_token')

	useEffect(() => {
		if (!userId || !accessToken) {
			message.warning('请先登录')
			navigate('/login')
			return
		}
		fetchProfile()
		fetchTransactions()
	}, [userId, accessToken])

	const fetchProfile = async () => {
		try {
			const response = await axios.get('/api/v1/user/profile', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setProfile(response.data.data)
				editForm.setFieldsValue({
					nickname: response.data.data.nickname,
					email: response.data.data.email,
					phone: response.data.data.phone
				})
			}
		} catch (error) {
			message.error('获取用户信息失败')
		} finally {
			setLoading(false)
		}
	}

	const fetchTransactions = async () => {
		try {
			const response = await axios.get('/api/v1/wallet/transactions', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setTransactions(response.data.data)
			}
		} catch (error) {
			message.error('获取交易记录失败')
		}
	}

	const handleUpdateProfile = async (values: any) => {
		try {
			const response = await axios.put('/api/v1/user/profile', values, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				message.success('更新成功')
				fetchProfile()
			} else {
				message.error(response.data.message || '更新失败')
			}
		} catch (error) {
			message.error('更新失败')
		}
	}

	const handleRecharge = async (values: { amount: number }) => {
		try {
			const response = await axios.post('/api/v1/wallet/recharge', values, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				message.success('充值成功')
				rechargeForm.resetFields()
				fetchProfile()
			} else {
				message.error(response.data.message || '充值失败')
			}
		} catch (error) {
			message.error('充值失败')
		}
	}

	const handleLogout = () => {
		localStorage.removeItem('access_token')
		localStorage.removeItem('refresh_token')
		localStorage.removeItem('user_id')
		localStorage.removeItem('nickname')
		message.success('已退出登录')
		navigate('/login')
	}

	if (loading) {
		return (
			<div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
				<Spin size="large" />
			</div>
		)
	}

	return (
		<div style={{ maxWidth: 800, margin: '0 auto', padding: '20px' }}>
			<Card title="用户设置">
				<Tabs
					activeKey={activeTab}
					onChange={setActiveTab}
					items={[
						{
							key: 'profile',
							label: '个人资料',
							children: (
								<Form
									form={editForm}
									layout="vertical"
									onFinish={handleUpdateProfile}
								>
									<div style={{ textAlign: 'center', marginBottom: 24 }}>
										<Avatar size={100} src={profile?.avatar_url} icon="user" />
										<Upload showUploadList={false}>
											<Button style={{ marginTop: 12 }}>更换头像</Button>
										</Upload>
									</div>

									<Form.Item label="用户名">
										<Input value={profile?.username} disabled />
									</Form.Item>

									<Form.Item
										name="nickname"
										label="昵称"
										rules={[{ required: true, message: '请输入昵称' }]}
									>
										<Input placeholder="请输入昵称" />
									</Form.Item>

									<Form.Item
										name="email"
										label="邮箱"
										rules={[{ type: 'email', message: '请输入有效邮箱' }]}
									>
										<Input placeholder="请输入邮箱" />
									</Form.Item>

									<Form.Item
										name="phone"
										label="手机号"
									>
										<Input placeholder="请输入手机号" />
									</Form.Item>

									<Form.Item>
										<Button type="primary" htmlType="submit" block>
											保存修改
										</Button>
									</Form.Item>
								</Form>
							)
						},
						{
							key: 'wallet',
							label: '钱包',
							children: (
								<div>
									<Card size="small" style={{ marginBottom: 16 }}>
										<div style={{ fontSize: 24, fontWeight: 'bold', color: '#ff6b00' }}>
											{profile?.coin_balance || 0} 虎牙币
										</div>
										<div style={{ color: '#999' }}>当前余额</div>
									</Card>

									<Form
										form={rechargeForm}
										layout="vertical"
										onFinish={handleRecharge}
									>
										<Form.Item
											name="amount"
											label="充值金额"
											rules={[{ required: true, message: '请输入充值金额' }]}
										>
											<Input type="number" min={1} placeholder="请输入充值金额" />
										</Form.Item>

										<Form.Item>
											<Button type="primary" htmlType="submit" block>
												充值
											</Button>
										</Form.Item>
									</Form>

									<div style={{ marginTop: 16 }}>
										<h4>交易记录</h4>
										<List
											size="small"
											dataSource={transactions.slice(0, 10)}
											renderItem={(item) => (
												<List.Item>
													<List.Item.Meta
														title={item.description}
														description={new Date(item.created_at).toLocaleString()}
													/>
													<span style={{ color: item.amount > 0 ? 'green' : 'red' }}>
														{item.amount > 0 ? '+' : ''}{item.amount}
													</span>
												</List.Item>
											)}
										/>
									</div>
								</div>
							)
						},
						{
							key: 'level',
							label: '等级',
							children: (
								<div style={{ textAlign: 'center', padding: 24 }}>
									<Avatar size={80} style={{ backgroundColor: '#ff6b00' }}>
										Lv.{profile?.level || 1}
									</Avatar>
									<h3 style={{ marginTop: 16 }}>当前等级: Lv.{profile?.level || 1}</h3>
									<p style={{ color: '#999' }}>
										经验值: {profile?.exp || 0}
									</p>
									<div style={{ marginTop: 24, textAlign: 'left' }}>
										<h4>升级攻略</h4>
										<ul style={{ color: '#666' }}>
											<li>观看直播可获得经验值</li>
											<li>发送弹幕可获得经验值</li>
											<li>赠送礼物可获得大量经验值</li>
											<li>开通会员可获得额外经验值</li>
										</ul>
									</div>
								</div>
							)
						},
						{
							key: 'security',
							label: '安全',
							children: (
								<div>
									<Button danger block onClick={handleLogout} style={{ marginBottom: 16 }}>
										退出登录
									</Button>
									<Card size="small" title="账号安全">
										<p style={{ color: '#666' }}>当前账号状态: 正常</p>
										<p style={{ color: '#666' }}>登录设备: Mac / Chrome</p>
									</Card>
								</div>
							)
						}
					]}
				/>
			</Card>
		</div>
	)
}

export default Settings
