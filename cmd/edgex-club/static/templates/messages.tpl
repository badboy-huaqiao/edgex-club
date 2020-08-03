{{define "messagesList"}}
{{range $key,$msg := .}}
<div  class="list-group-item list-group-item-action">
        <a href="#" onclick="updateMsgAndSkip({{$msg.ArticleUserName}},{{$msg.ArticleId}},{{$msg.Id}})">{{$msg.Content}}</a>
        <span class="pull-right">
            {{if $msg.IsRead }}
            <span class="badge badge-pill badge-secondary">已读</span> 
            {{else}} 
            <span class="badge badge-pill badge-danger">未读</span>
            {{end}}
            <span class="badge badge-pill badge-success">{{$msg.Created | fdate}}</span>
        </span>
        
</div> 
{{end}}
{{end}}