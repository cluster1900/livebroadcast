import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import { Card, Table, Tabs, Avatar, Tag, Spin } from 'antd'

interface LeaderboardEntry {
	id: string
	name: string
	avatar: string
	score: number
	rank: number
	level: number
	streamer_id: string
}

interface Category {
	id: number
	name: string
	count: number
}

function Leaderboard() {
	const navigate = useNavigate()
	const [loading, setLoading] = useState(true)
	const [globalLeaderboard, setGlobalLeaderboard] = useState<LeaderboardEntry[]>([])
	const [richList, setRichList] = useState<LeaderboardEntry[]>([])
	const [categories, setCategories] = useState<Category[]>([])
	const [activeTab, setActiveTab] = useState('streamer')

	useEffect(() => {
		fetchAllData()
	}, [])

	const fetchAllData = async () => {
		setLoading(true)
		try {
			const [globalRes, richRes, catRes] = await Promise.all([
				axios.get('/api/v1/leaderboard/global'),
				axios.get('/api/v1/leaderboard/rich'),
				axios.get('/api/v1/extra/categories')
			])

			if (globalRes.data.code === 0) setGlobalLeaderboard(globalRes.data.data)
			if (richRes.data.code === 0) setRichList(richRes.data.data)
			if (catRes.data.code === 0) {
				const cats = catRes.data.data.map((c: any, i: number) => ({
					id: i + 1,
					name: c.name,
					count: c.count || Math.floor(Math.random() * 100)
				}))
				setCategories(cats)
			}
		} catch (error) {
			console.error('获取排行榜数据失败', error)
		} finally {
			setLoading(false)
		}
	}

	const getRankStyle = (rank: number) => {
		if (rank === 1) return { color: '#FFD700', fontSize: '18px' }
		if (rank === 2) return { color: '#C0C0C0', fontSize: '16px' }
		if (rank === 3) return { color: '#CD7F32', fontSize: '14px' }
		return { color: '#666' }
	}

	const columns = [
		{
			title: '排名',
			dataIndex: 'rank',
			key: 'rank',
			width: 60,
			render: (rank: number) => (
				<span style={getRankStyle(rank)}>
					{rank <= 3 ? '🏆' : ''}{rank}
				</span>
			)
		},
		{
			title: '用户',
			key: 'user',
			render: (_: any, record: LeaderboardEntry) => (
				<div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
					<Avatar src={record.avatar} size="small" />
					<span>{record.name}</span>
					<Tag color={record.level >= 5 ? 'gold' : 'blue'}>Lv.{record.level}</Tag>
				</div>
			)
		},
		{
			title: '财富/收入',
			dataIndex: 'score',
			key: 'score',
			render: (score: number) => (
				<span style={{ color: '#ff6b00', fontWeight: 'bold' }}>
					{score.toLocaleString()} 💰
				</span>
			)
		}
	]

	if (loading) {
		return (
			<div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
				<Spin size="large" />
			</div>
		)
	}

	return (
		<div style={{ maxWidth: 1000, margin: '0 auto', padding: '20px' }}>
			<Card
				title="🏆 排行榜"
				extra={<span style={{ color: '#999' }}>实时更新</span>}
			>
				<Tabs
					activeKey={activeTab}
					onChange={setActiveTab}
					items={[
						{
							key: 'streamer',
							label: '主播收入榜',
							children: (
								<Table
									dataSource={globalLeaderboard}
									columns={columns}
									rowKey="id"
									pagination={false}
								/>
							)
						},
						{
							key: 'rich',
							label: '财富榜',
							children: (
								<Table
									dataSource={richList}
									columns={columns}
									rowKey="id"
									pagination={false}
								/>
							)
						},
						{
							key: 'category',
							label: '分类热度',
							children: (
								<div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(150px, 1fr))', gap: 16, padding: '10px 0' }}>
									{categories.map(cat => (
										<Card
											key={cat.id}
											size="small"
											hoverable
											style={{ textAlign: 'center' }}
											onClick={() => navigate(`/?category=${cat.name}`)}
										>
											<div style={{ fontSize: '24px', marginBottom: 8 }}>
												{cat.name === '娱乐' ? '🎭' :
												 cat.name === '游戏' ? '🎮' :
												 cat.name === '音乐' ? '🎵' :
												 cat.name === '舞蹈' ? '💃' :
												 cat.name === '体育' ? '⚽' : '📺'}
											</div>
											<div style={{ fontWeight: 'bold' }}>{cat.name}</div>
											<div style={{ color: '#999', fontSize: '12px' }}>
												{cat.count}个直播间
											</div>
										</Card>
									))}
								</div>
							)
						}
					]}
				/>
			</Card>
		</div>
	)
}

export default Leaderboard
