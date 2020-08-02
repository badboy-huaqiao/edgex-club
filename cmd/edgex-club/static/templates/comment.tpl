
{{define "comment"}}
<div class="list-group-item list-group-item-action">
    <div class="media">
        <a href="/user/{{.UserName}}"><img src="{{.UserAvatarUrl}}" class="mr-3 rounded-circle" style="width:45px;height:45px; " alt="..."></a>
        <div class="media-body">
            <p class="mt-0 mb-1">{{.Content}}</p>
            <span class="badge badge-secondary">{{ fdate .Created }}</span>
            <span class="badge badge-secondary pull-right" onclick="showReply({{.Id.Hex}})">回复</span>
        </div>
    </div>
    <div id="{{.Id.Hex}}" style="display: none;margin-top:5px;" class="list-group-item list-group-item-action list-group-item-secondary">
        <textarea class="reply_context" style="outline:none;" placeholder="回复@：{{.UserName}}"></textarea>
        <span class="badge badge-secondary pull-right" onclick="reply({{.Id.Hex}},{{.UserName}},{{.Id.Hex}})">提交回复</span>
    </div>
    <div class="list-group list-group-flush" id="{{.Id.Hex}}-reply-list">
    </div>   
</div>
{{end}}