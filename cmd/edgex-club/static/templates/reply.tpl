 {{define "reply"}}
<div class="list-group-item list-group-item-action list-group-item-secondary">
    <div class="media">
        <div class="media-body">
            <p class="mt-0">
                <a href="/user/{{.FromUserName}}">{{.FromUserName}}</a>：回复 <a href="/user/{{.ToUserName}}">@{{.ToUserName}}：</a> {{.Content}}
            </p>
            <span class="badge badge-secondary">{{.Created | fdate}}</span>
            <span class="badge badge-secondary pull-right" onclick="showReply({{.Id.Hex}})">回复</span>
        </div>
    </div>
    <div id="{{.Id.Hex}}" style="display: none;margin-top:5px;">
        <textarea class="reply_context" style="outline:none;" placeholder="回复@：{{.FromUserName}}"></textarea>
        <span class="badge badge-secondary pull-right" onclick="reply({{.CommentId}},{{.FromUserName}},{{.Id.Hex}})">提交回复</span>
    </div>
</div>
{{end}}
