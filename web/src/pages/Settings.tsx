import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import { message, Card, Form, Input, Button, Upload, Avatar, Tabs, List, Spin, Progress, Modal } from 'antd'
import { UserOutlined, LockOutlined, UploadOutlined } from '@ant-design/icons'

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
	const [passwordForm] = Form.useForm()
	const [rechargeForm] = Form.useForm()
	const [activeTab, setActiveTab] = useState('profile')
	const [uploading, setUploading] = useState(false)
	const [passwordModalOpen, setPasswordModalOpen] = useState(false)

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

	const handleAvatarUpload = async (file: File) => {
		setUploading(true)
		try {
			const formData = new FormData()
			formData.append('file', file)

			const response = await axios.post('/api/v1/upload/avatar', formData, {
				headers: {
					Authorization: `Bearer ${accessToken}`,
					'Content-Type': 'multipart/form-data'
				}
			})

			if (response.data.code === 0) {
				const avatarUrl = response.data.data.url
				await axios.put('/api/v1/user/profile', { avatar_url: avatarUrl }, {
					headers: { Authorization: `Bearer ${accessToken}` }
				})
				message.success('头像上传成功')
				fetchProfile()
			} else {
				message.error(response.data.message || '上传失败')
			}
		} catch (error) {
			message.error('头像上传失败，请重试')
		} finally {
			setUploading(false)
		}
	}

	const handleChangePassword = async (values: { old_password: string; new_password: string }) => {
		try {
			const response = await axios.post('/api/v1/password/change', values, {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				message.success('密码修改成功')
				setPasswordModalOpen(false)
				passwordForm.resetFields()
			} else {
				message.error(response.data.message || '修改失败')
			}
		} catch (error: any) {
			message.error(error.response?.data?.message || '密码修改失败')
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
		Modal.confirm({
			title: '确认退出',
			content: '确定要退出登录吗？',
			onOk: () => {
				localStorage.removeItem('access_token')
				localStorage.removeItem('refresh_token')
				localStorage.removeItem('user_id')
				localStorage.removeItem('nickname')
				message.success('已退出登录')
				navigate('/login')
			}
		})
	}

	const getLevelProgress = (level: number, exp: number) => {
		const levelExp = [0, 2000, 5000, 10000, 20000, 50000, 100000, 200000, 500000, 1000000]
		const nextExp = levelExp[level] || level * 500000
		const prevExp = levelExp[level - 1] || 0
		if (exp >= nextExp) return 100
		return Math.round(((exp - prevExp) / (nextExp - prevExp)) * 100)
	}

	if (loading) {
		return (
			<div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
				<Spin size="large" />
			</div>
		)
	}

	return (
		<div style={{ maxWidth: 900, margin: '0 auto', padding: '20px' }}>
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
										<Upload
											showUploadList={false}
											beforeUpload={(file) => {
												handleAvatarUpload(file)
												return false
											}}
										>
											<div style={{ cursor: 'pointer' }}>
												<Avatar size={120} src={profile?.avatar_url} icon={<UserOutlined />} style={{ backgroundColor: '#ff6b00' }} />
												{uploading && <Progress percent={50} size="small" style={{ marginTop: 8 }} />}
												<div style={{ marginTop: 8, color: '#1890ff' }}>
													<UploadOutlined /> 点击更换头像
												</div>
											</div>
										</Upload>
									</div>

									<Form.Item label="用户名">
										<Input value={profile?.username} disabled prefix={<UserOutlined />} />
									</Form.Item>

									<Form.Item
										name="nickname"
										label="昵称"
										rules={[{ required: true, message: '请输入昵称', min: 2, max: 20 }]}
									>
										<Input placeholder="请输入昵称（2-20个字符）" />
									</Form.Item>

									<Form.Item
										name="email"
										label="邮箱"
										rules={[{ type: 'email', message: '请输入有效邮箱' }]}
									>
										<Input placeholder="请输入邮箱" />
									</Form.Item>

									<Form.Item name="phone" label="手机号">
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
							key: 'security',
							label: '账号安全',
							children: (
								<div>
									<Card size="small" title="密码管理" style={{ marginBottom: 16 }}>
										<p style={{ color: '#666', marginBottom: 16 }}>
											建议定期更换密码，确保账号安全
										</p>
										<Button
											type="primary"
											icon={<LockOutlined />}
											onClick={() => setPasswordModalOpen(true)}
										>
											修改密码
										</Button>
									</Card>

									<Card size="small" title="设备管理" style={{ marginBottom: 16 }}>
										<p style={{ color: '#666' }}>当前登录设备: Mac / Chrome</p>
										<p style={{ color: '#999', fontSize: 12 }}>IP: 192.168.1.1</p>
									</Card>

									<Card size="small" title="绑定信息">
										<p style={{ color: '#666' }}>
											邮箱: {profile?.email || '未绑定'}
										</p>
										<p style={{ color: '#666' }}>
											手机: {profile?.phone || '未绑定'}
										</p>
									</Card>

									<Button danger block onClick={handleLogout} style={{ marginTop: 16 }}>
										退出登录
									</Button>
								</div>
							)
						},
						{
							key: 'wallet',
							label: '钱包',
							children: (
								<div>
									<Card size="small" style={{ marginBottom: 16, background: 'linear-gradient(135deg, #ff6b00 0%, #ff9a44 100%)' }}>
										<div style={{ fontSize: 28, fontWeight: 'bold', color: '#fff' }}>
											{profile?.coin_balance || 0} 虎牙币
										</div>
										<div style={{ color: 'rgba(255,255,255,0.8)' }}>当前余额</div>
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

									<div style={{ marginTop: 24 }}>
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
													<span style={{ color: item.amount > 0 ? 'green' : 'red', fontWeight: 'bold' }}>
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
							label: '等级特权',
							children: (
								<div style={{ padding: 16 }}>
									<div style={{ textAlign: 'center', marginBottom: 24 }}>
										<Avatar size={100} style={{ backgroundColor: '#ff6b00', fontSize: 36 }}>
											Lv.{profile?.level || 1}
										</Avatar>
										<h3 style={{ marginTop: 12 }}>当前等级: Lv.{profile?.level || 1}</h3>
										<p style={{ color: '#999' }}>
											经验值: {profile?.exp || 0}
										</p>
										<Progress
											percent={getLevelProgress(profile?.level || 1, profile?.exp || 0)}
											strokeColor="#ff6b00"
											style={{ maxWidth: 300, margin: '0 auto' }}
										/>
									</div>

									<Card size="small" title="等级特权" style={{ marginBottom: 16 }}>
										<List size="small">
											<List.Item>弹幕特权 - 更高频率发送弹幕</List.Item>
											<List.Item>礼物特权 - 专属礼物和动画</List.Item>
											<List.Item>等级标识 - 弹幕显示等级</List.Item>
											<List.Item>优先客服 - 优先获得客服支持</List.Item>
										</List>
									</Card>

									<Card size="small" title="升级攻略">
										<List size="small">
											<List.Item>观看直播 1分钟 = 1经验值</List.Item>
											<List.Item>发送弹幕 1条 = 10经验值</List.Item>
											<List.Item>赠送礼物 1虎牙币 = 1经验值</List.Item>
											<List.Item>每日首次登录 = 10经验值</List.Item>
										</List>
									</Card>
								</div>
							)
						}
					]}
				/>
			</Card>

			<Modal
				title="修改密码"
				open={passwordModalOpen}
				onCancel={() => setPasswordModalOpen(false)}
				footer={null}
			>
				<Form form={passwordForm} layout="vertical" onFinish={handleChangePassword}>
					<Form.Item
						name="old_password"
						label="原密码"
						rules={[{ required: true, message: '请输入原密码' }]}
					>
						<Input.Password prefix={<LockOutlined />} placeholder="请输入原密码" />
					</Form.Item>
					<Form.Item
						name="new_password"
						label="新密码"
						rules={[
							{ required: true, message: '请输入新密码' },
							{ min: 6, message: '密码长度至少6位' }
						]}
					>
						<Input.Password prefix={<LockOutlined />} placeholder="请输入新密码（至少6位）" />
					</Form.Item>
					<Form.Item
						name="confirm_password"
						label="确认新密码"
						dependencies={['new_password']}
						rules={[
							{ required: true, message: '请确认新密码' },
							({ getFieldValue }) => ({
								validator(_, value) {
									if (!value || getFieldValue('new_password') === value) {
										return Promise.resolve()
									}
									return Promise.reject(new Error('两次输入的密码不一致'))
								}
							})
						]}
					>
						<Input.Password prefix={<LockOutlined />} placeholder="请再次输入新密码" />
					</Form.Item>
					<Form.Item>
						<Button type="primary" htmlType="submit" block>
							确认修改
						</Button>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	)
}

export default Settings
