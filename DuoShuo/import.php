<?php

// 读取 json 文件并转换成 php 数组
$json = file_get_contents("./export.json");
$data = json_decode($json, true);

// 文章数据
$threads = $data['threads'];	
// 评论数据
$posts = $data['posts'];	

// 多说的文章ID(thread_id)与 typecho 的文章ID(cid) 对应关系
$threadIdRelationCid = [];
foreach ($threads as $item) {
	$threadIdRelationCid[$item['thread_id']] = $item['thread_key'];
}

// 假设评论表 coid 小于10001，此处从1001开始自增，请根据实际最大值修改
$coid = 10001;
// 多说的评论ID(post_id) 与 typecho 的评论ID(coid) 对应关系
$postIdRelationCoid = [];
foreach ($posts as $item) {
	$postIdRelationCoid[$item['post_id']] = $coid++;
}

// 拼成多条 insert sql语句
$sql = '';
foreach ($posts as $item) {
	$coid = $postIdRelationCoid[$item['post_id']];
	$cid = $threadIdRelationCid[$item['thread_id']];
	$created = strtotime($item['created_at']);
	$author = $item['author_name'] ?: '';
	$mail = $item['author_email'] ?: '';
	$url = $item['author_url'] ?: '';
	$ip = $item['ip'];
	$text = $item['message'];
	$parent = 0;
	if (is_array($item['parents'])) {
		$parent = $postIdRelationCoid[$item['parents'][0]];
	}

	$sql .= "INSERT INTO `typecho_comments` 
(`coid`, `cid`, `created`, `author`, `authorId`, `ownerId`, `mail`, `url`, `ip`, `agent`, `text`, `type`, `status`, `parent`) VALUES
({$coid}, {$cid}, {$created}, '{$author}', 0, 1, '{$mail}', '{$url}', '{$ip}', NULL, '{$text}', 'comment', 'approved', $parent);\n";
}

// 将 sql 写入文件中
file_put_contents("./insert.sql", $sql);

echo "end \n";

/*
INSERT INTO `typecho_comments` 
(`coid`, `cid`, `created`, `author`, `authorId`, `ownerId`, `mail`, `url`, `ip`, `agent`, `text`, `type`, `status`, `parent`) VALUES
(8, 402, 1493430000, 'sean', 0, 1, 'json_vip@163.com', NULL, '::1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36', '赞的回复', 'comment', 'approved', 6);
*/















