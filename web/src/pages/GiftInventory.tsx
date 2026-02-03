import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import { message, Card, Spin, Empty, Tag } from 'antd'
import { GiftOutlined } from '@ant-design/icons'

interface GiftItem {
	gift_id: number
	name: string
	icon_url: string
	count: number
	coin_price: number
	category: string
}

function GiftInventory() {
	const navigate = useNavigate()
	const [loading, setLoading] = useState(true)
	const [gifts, setGifts] = useState<GiftItem[]>([])

	const accessToken = localStorage.getItem('access_token')

	useEffect(() => {
		if (!accessToken) {
			message.warning('请先登录')
			navigate('/login')
			return
		}
		fetchGifts()
	}, [accessToken])

	const fetchGifts = async () => {
		try {
			const response = await axios.get('/api/v1/inventory/gifts', {
				headers: { Authorization: `Bearer ${accessToken}` }
			})
			if (response.data.code === 0) {
				setGifts(response.data.data)
			}
		} catch (error) {
			message.error('获取礼物背包失败')
		} finally {
			setLoading(false)
		}
	}

	const getCategoryColor = (category: string) => {
		switch (category) {
			case 'vip': return 'gold'
			case 'special': return 'purple'
			default: return 'blue'
		}
	}

	const totalValue = gifts.reduce((sum, g) => sum + g.count * g.coin_price, 0)

	if (loading) {
		return (
			<div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
				<Spin size="large" />
			</div>
		)
	}

	return (
		<div style={{ maxWidth: 1200, margin: '0 auto', padding: '20px' }}>
			<Card
				title={
					<span>
						<GiftOutlined style={{ marginRight: 8 }} />
						礼物背包
					</span>
				}
				extra={
					<Tag color="orange">总价值: {totalValue} 虎牙币</Tag>
				}
			>
				{gifts.length === 0 ? (
					<Empty
						description="礼物背包为空"
						style={{ padding: 60 }}
					>
						<div style={{ marginBottom: 16 }}>
							去直播间送礼物可以获得更多礼物哦！
						</div>
						<a onClick={() => navigate('/')}>去逛逛</a>
					</Empty>
				) : (
					<div style={{
						display: 'grid',
						gridTemplateColumns: 'repeat(auto-fill, minmax(150px, 1fr))',
						gap: 16
					}}>
						{gifts.map((gift) => (
							<Card
								key={gift.gift_id}
								size="small"
								hoverable
								style={{ textAlign: 'center' }}
							>
								<div style={{ fontSize: 48, marginBottom: 8 }}>
									{gift.icon_url}
								</div>
								<div style={{ fontWeight: 'bold', marginBottom: 4 }}>
									{gift.name}
								</div>
								<div style={{ color: '#ff6b00', fontSize: 14, marginBottom: 4 }}>
									{gift.coin_price} 虎牙币
								</div>
								<div style={{
									background: '#f5f5f5',
									borderRadius: 20,
									padding: '4px 12px',
									display: 'inline-block'
								}}>
									x{gift.count}
								</div>
								<div style={{ marginTop: 8 }}>
									<Tag color={getCategoryColor(gift.category)}>
										{gift.category === 'vip' ? 'VIP专属' : gift.category === 'special' ? 'Special' : '普通'}
									</Tag>
								</div>
							</Card>
						))}
					</div>
				)}
			</Card>
		</div>
	)
}

export default GiftInventory
